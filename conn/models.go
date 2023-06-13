/*
thanks for the base structs mtarnawa/godesu, even though
they were shit.
*/

package conn

import (
	"net/url"
)

type Image struct {
	URL                url.URL
	Filename, Original string
}

type Post struct {
	IsDeleted bool  `json:"is_deleted"`
	Parent    int   `json:"parent"`
	Image     Image `json:"image"`

	No          int    `json:"no"`
	Sticky      int    `json:"sticky"`
	Closed      int    `json:"closed"`
	Now         string `json:"now"`
	Name        string `json:"name"`
	Sub         string `json:"sub,omitempty"`
	Com         string `json:"com"`
	Filename    string `json:"filename,omitempty"`
	Ext         string `json:"ext,omitempty"`
	W           int    `json:"w,omitempty"`
	H           int    `json:"h,omitempty"`
	TnW         int    `json:"tn_w,omitempty"`
	TnH         int    `json:"tn_h,omitempty"`
	Tim         int64  `json:"tim,omitempty"`
	Time        int64  `json:"time"`
	Md5         string `json:"md5,omitempty"`
	Fsize       int    `json:"fsize,omitempty"`
	Resto       int    `json:"resto"`
	Bumplimit   int    `json:"bumplimit,omitempty"`
	Imagelimit  int    `json:"imagelimit,omitempty"`
	SemanticURL string `json:"semantic_url,omitempty"`
	Replies     int    `json:"replies,omitempty"`
	Images      int    `json:"images,omitempty"`
	UniqueIps   int    `json:"unique_ips,omitempty"`
	Trip        string `json:"trip,omitempty"`
}

func (p Post) IsOP() bool {
	return p.Resto == 0
}

type SmallThread struct {
	IsDeleted bool `json:"is_deleted"`
	Page      int  `json:"current_page"`

	No           int `json:"no"`
	LastModified int `json:"last_modified"`
	Replies      int `json:"replies"`
	LastReplies  int `json:"last_replies"`
}

type SmallPage struct {
	Threads []SmallThread `json:"threads"`
	Page    int           `json:"page"`
}

type Buf struct {
	TotalThreads int

	Posts   []Post      `json:"posts"`
	Page    int         `json:"page"`
	Threads []SmallPage `json:"threads"`

	//No           int `json:"no"`
	//LastModified int `json:"last_modified"`
	//Replies      int `json:"replies"`
}

/* returns -1 if not present */
func (b Buf) Locate(no int) int {
	for k, v := range b.Posts {
		if v.No == no {
			return k
		}
	}
	return -1
}

/* returns (-1, -1) if not present */
func (b Buf) LocateThread(no int) (int, int) {
	for k, v := range b.Threads {
		for kk, vv := range v.Threads {
			if vv.No == no {
				return k, kk
			}
		}
	}
	return -1, -1
}

type IntrBoard struct {
	Posts       *IntPost
	Threads     *IntST
	JsonPosts   map[int]Post        `json:"map_Posts"`
	JsonThreads map[int]SmallThread `json:"map_SmallThreads"`
}

type Intr struct {
	Active  IntrBoard
	Deleted IntrBoard
}

/* one type to kill em all? */
type Bored struct {
	Name string

	Intr Intr
	Buf  Buf /* doesn't need to be saved between runs */

	Worker Worker
}

func mkBored(name string) Bored {
	return Bored{
		Name: name,
		Intr: Intr{
			Active: IntrBoard{
				Posts:   NewEmptyIntPost(),
				Threads: NewEmptyIntST(),
			},
			Deleted: IntrBoard{
				Posts:   NewEmptyIntPost(),
				Threads: NewEmptyIntST(),
			},
		},
		Buf: Buf{},
	}
}

/* 4LS tries to only use above structs */

type Board struct {
	SmallPages   []SmallPage /* buffer */
	SmallThreads map[int]*SmallThread

	Board           string    `json:"board"`
	Title           string    `json:"title"`
	WsBoard         int       `json:"ws_board"`
	PerPage         int       `json:"per_page"`
	PagesN          int       `json:"pages"`
	MaxFilesize     int       `json:"max_filesize"`
	MaxWebmFilesize int       `json:"max_webm_filesize"`
	MaxCommentChars int       `json:"max_comment_chars"`
	MaxWebmDuration int       `json:"max_webm_duration"`
	BumpLimit       int       `json:"bump_limit"`
	ImageLimit      int       `json:"image_limit"`
	Cooldowns       Cooldowns `json:"cooldowns"`
	MetaDescription string    `json:"meta_description"`
	IsArchived      int       `json:"is_archived,omitempty"`
	Spoilers        int       `json:"spoilers,omitempty"`
	CustomSpoilers  int       `json:"custom_spoilers,omitempty"`
	ForcedAnon      int       `json:"forced_anon,omitempty"`
	UserIds         int       `json:"user_ids,omitempty"`
	CountryFlags    int       `json:"country_flags,omitempty"`
	CodeTags        int       `json:"code_tags,omitempty"`
	WebmAudio       int       `json:"webm_audio,omitempty"`
	MinImageWidth   int       `json:"min_image_width,omitempty"`
	MinImageHeight  int       `json:"min_image_height,omitempty"`
	Oekaki          int       `json:"oekaki,omitempty"`
	SjisTags        int       `json:"sjis_tags,omitempty"`
	TextOnly        int       `json:"text_only,omitempty"`
	RequireSubject  int       `json:"require_subject,omitempty"`
	TrollFlags      int       `json:"troll_flags,omitempty"`
	MathTags        int       `json:"math_tags,omitempty"`

	Threads map[int]Thread /* not used */
	Buf     []CataPage     /* not used */
	Pages   []Page         /* not used */
}

/* only used to unmarshal thread requests */
type Posts struct {
	P []Post `json:"posts"`
}

type LastReplies struct {
	No       int    `json:"no"`
	Now      string `json:"now"`
	Name     string `json:"name"`
	Com      string `json:"com"`
	Filename string `json:"filename,omitempty"`
	Ext      string `json:"ext,omitempty"`
	W        int    `json:"w,omitempty"`
	H        int    `json:"h,omitempty"`
	TnW      int    `json:"tn_w,omitempty"`
	TnH      int    `json:"tn_h,omitempty"`
	Tim      int64  `json:"tim,omitempty"`
	Time     int    `json:"time"`
	Md5      string `json:"md5,omitempty"`
	Fsize    int    `json:"fsize,omitempty"`
	Resto    int    `json:"resto"`
}

type Thread struct {
	No            int           `json:"no"`
	Sticky        int           `json:"sticky,omitempty"`
	Closed        int           `json:"closed,omitempty"`
	Now           string        `json:"now"`
	Name          string        `json:"name,omitempty"`
	Sub           string        `json:"sub,omitempty"`
	Com           string        `json:"com"`
	Filename      string        `json:"filename"`
	Ext           string        `json:"ext"`
	W             int           `json:"w"`
	H             int           `json:"h"`
	TnW           int           `json:"tn_w"`
	TnH           int           `json:"tn_h"`
	Tim           int64         `json:"tim"`
	Time          int           `json:"time"`
	Md5           string        `json:"md5"`
	Fsize         int           `json:"fsize"`
	Resto         int           `json:"resto"`
	SemanticURL   string        `json:"semantic_url"`
	Replies       int           `json:"replies"`
	Images        int           `json:"images"`
	LastModified  int           `json:"last_modified"`
	Bumplimit     int           `json:"bumplimit,omitempty"`
	Imagelimit    int           `json:"imagelimit,omitempty"`
	OmittedPosts  int           `json:"omitted_posts,omitempty"`
	OmittedImages int           `json:"omitted_images,omitempty"`
	LastReplies   []LastReplies `json:"last_replies,omitempty"`
	Trip          string        `json:"trip,omitempty"`
	BufPosts      Posts         `json:"posts"`

	Posts map[int]Post
}

/* unclean from godesu */
type Page struct {
	board string
	// c     *Gochan
	All []struct {
		Posts []struct {
			No          int    `json:"no"`
			Now         string `json:"now"`
			Name        string `json:"name"`
			Sub         string `json:"sub"`
			Com         string `json:"com"`
			Filename    string `json:"filename"`
			Ext         string `json:"ext"`
			W           int    `json:"w"`
			H           int    `json:"h"`
			TnW         int    `json:"tn_w"`
			TnH         int    `json:"tn_h"`
			Tim         int64  `json:"tim"`
			Time        int    `json:"time"`
			Md5         string `json:"md5"`
			Fsize       int    `json:"fsize"`
			Resto       int    `json:"resto"`
			SemanticURL string `json:"semantic_url"`
			Replies     int    `json:"replies"`
			Images      int    `json:"images"`
		} `json:"posts"`
	} `json:"threads"`
}

type Cooldowns struct {
	Threads int `json:"threads"`
	Replies int `json:"replies"`
	Images  int `json:"images"`
}

type CataPage struct {
	Page    int      `json:"page"`
	Threads []Thread `json:"threads"`
}

type SiteInfo struct {
	All        []Board `json:"boards"`
	TrollFlags struct {
		AC string `json:"AC"`
		AN string `json:"AN"`
		BL string `json:"BL"`
		CF string `json:"CF"`
		CM string `json:"CM"`
		CT string `json:"CT"`
		DM string `json:"DM"`
		EU string `json:"EU"`
		FC string `json:"FC"`
		GN string `json:"GN"`
		GY string `json:"GY"`
		JH string `json:"JH"`
		KN string `json:"KN"`
		MF string `json:"MF"`
		NB string `json:"NB"`
		NZ string `json:"NZ"`
		PC string `json:"PC"`
		PR string `json:"PR"`
		RE string `json:"RE"`
		TM string `json:"TM"`
		TR string `json:"TR"`
		UN string `json:"UN"`
		WP string `json:"WP"`
	} `json:"troll_flags"`
}
