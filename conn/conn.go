package conn

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	errs "github.com/pkg/errors"

	// "github.com/lordrusk/4LS/bord"
	"github.com/lordrusk/4LS/conn/logger"
	// fson "github.com/multiprocessio/go-json"
)

var (
	scheme string = "https"

	/* hosts */
	Static     string = "a.4cdn.org"
	Img        string = "i.4cdn.org"
	StaticInfo string = "s.4cdn.org" /* probably unused, kept for expansion */
)

type Con struct {
	*http.Client
	logger.Logger
	W *Worker /* emptry Worker.work func */
	*allbrds
	lp string /* log folder */
}

func NewCon(logFolder string) *Con {
	c := &Con{
		Client: &http.Client{
			Timeout: time.Second * 60,
		},
		Logger: logger.New(os.Stdout, "", log.Ltime|log.Ldate,
			logFolder+"/logs/4LS_"+time.Now().Format("2006-01-02_15:04:05")+".log"),
		allbrds: mkAll(),
		lp:      logFolder,
	}
	w := NewCeo(func(_ *Con) {}, c, true)
	c.W = w
	return c
}

/* model must be pointer to interface
 *
 * for images use Con.(*http.Client.)Get */
func (c *Con) Gather(u url.URL, model interface{}) error {
	resp, err := c.Get(u.String()) /* implicit .json */
	if err != nil {
		return errs.Wrapf(err, "Failed to get %s", u.String())
	} else if resp.StatusCode != http.StatusOK {
		return errs.Errorf("http status not ok: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	bites, err := io.ReadAll(resp.Body)
	if err != nil {
		return errs.Wrap(err, "couldnt read body")
	}
	// c.Printf("%s\n", bites)
	return json.Unmarshal(bites, model)
}

/* model should be a pointer */
func GetModel(path string, model interface{}) error {
	path = strings.ReplaceAll(path, "~", "$HOME")
	path = os.ExpandEnv(path)
	bites, err := os.ReadFile(path)
	if err != nil {
		return errs.Wrapf(err, "Failed to read %s", path)
	}
	/* if err := json.Unmarshal(bytes.TrimSpace(bites), model); err != nil {/* */
	if err := json.Unmarshal(bites, model); err != nil { /* */
		return errs.Wrapf(err, "Failed to unmarshal %s into struct", path)
	}
	return nil
}

func PutModel(path string, model interface{}) error {
	path = strings.ReplaceAll(path, "~", "$HOME")
	path = os.ExpandEnv(path)
	spath := strings.Split(path, "/")

	var f *os.File
	var err error
	if f, err = os.OpenFile(path, os.O_RDWR, 0666); err != nil {
		if err := os.MkdirAll(strings.Join(spath[0:len(spath)-1], "/"), os.ModePerm); err != nil {
			return errs.Wrapf(err, "Cannot create '%s'", path)
		}
		if _, err = os.Create(path); err != nil {
			return errs.Wrapf(err, "Failed to create %s", path)
		}
		if f, err = os.OpenFile(path, os.O_RDWR, 0666); err != nil {
			return errs.Wrapf(err, "Failed to open %s", path)
		}
	}
	defer f.Close()

	/* bites, err := json.MarshalIndent(model, "", "	") /* */
	bites, err := json.Marshal(model) /* */
	if err != nil {
		return errs.Wrapf(err, "Failde to marshal struct")
	}

	/* almost instant json encodings */
	/* if err := fson.Encode(f, model); err != nil {
		return errs.Wrap(err, "Failed to marshal struct")
	} */
	if err := os.WriteFile(path, bites, 0666); err != nil {
		return errs.Wrapf(err, "Failed to write structure to %s", path)
	}

	return nil
}

func (c *Con) GetBoard(b *Bored) error {
	if err := GetModel(c.lp+"/"+b.Name+"_active.json", &b.Intr.Active); err != nil {
		return err
	}
	/* dont get previously deleted
	// if err := GetModel(c.lp+"/"+b.Name+"_deleted.json", &b.Intr.Deleted); err != nil {
	// 	return err
	// }
	*/
	b.Intr.Active.Posts.NewUnderlyer(b.Intr.Active.JsonPosts)
	return nil
}

/* puts deleted but doesnt get them, keeps writing over old deleted. Figure out how to single process that every save */
func (c *Con) PutBoard(b *Bored) error {
	// m := b.Intr.Active.Posts.GetMap()
	// b.Intr.Active.JsonPosts = m
	// b.Intr.Active.Posts.Clean()
	// mm := b.Intr.Active.Threads.GetMap()
	// b.Intr.Active.JsonThreads = mm
	// b.Intr.Active.Threads.Clean()
	// m = b.Intr.Deleted.Posts.GetMap()
	// b.Intr.Deleted.JsonPosts = m
	// b.Intr.Deleted.Posts.Clean()
	// mm = b.Intr.Deleted.Threads.GetMap()
	// b.Intr.Deleted.JsonThreads = mm
	// b.Intr.Deleted.Threads.Clean()

	m := b.Intr.Active.Posts.Close()
	b.Intr.Active.Posts = NewEmptyIntPost()
	b.Intr.Active.JsonPosts = m
	mm := b.Intr.Active.Threads.Close()
	b.Intr.Active.Threads = NewEmptyIntST()
	b.Intr.Active.JsonThreads = mm
	m = b.Intr.Deleted.Posts.Close()
	b.Intr.Deleted.Posts = NewEmptyIntPost()
	b.Intr.Deleted.JsonPosts = m
	mm = b.Intr.Deleted.Threads.Close()
	b.Intr.Deleted.Threads = NewEmptyIntST()
	b.Intr.Deleted.JsonThreads = mm

	if err := PutModel(c.lp+"/"+b.Name+"_active.json", &b.Intr.Active); err != nil {
		return err
	}
	if err := PutModel(c.lp+"/"+b.Name+"_deleted.json", &b.Intr.Deleted); err != nil {
		return err
	}

	b.Intr.Active.JsonPosts = make(map[int]Post)
	b.Intr.Active.JsonThreads = make(map[int]SmallThread)
	b.Intr.Deleted.JsonPosts = make(map[int]Post)
	b.Intr.Deleted.JsonThreads = make(map[int]SmallThread)
	return nil
}

/* url helper */
func Url(host, path string) url.URL {
	return url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}
}

/* helper function */
func MakePath(board string, number int) string {
	return fmt.Sprintf("/%s/thread/%d.json", board, number)
}

/* number is a string because of MakeImageN */
func MakeImgPath(extension, board, number string) string {
	return fmt.Sprintf("/%s/%s%s", board, number, extension)
}

/* get image number
 *
 * takes thread.Tim*/
func MakeImgN(tim int64) string {
	return strconv.FormatInt(tim, 10)
}

func MakeImgURL(endpoint string) url.URL {
	return url.URL{
		Scheme: scheme,
		Host:   Img,
		Path:   endpoint,
	}
}
