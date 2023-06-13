package main

import (
	// "fmt"

	errs "github.com/pkg/errors"

	"github.com/lordrusk/4LS/conn"
)

func deletedThreads(c *conn.Con, b *conn.Bored) {
	for stw := range b.Intr.Active.Threads.Range() {
		if !stw.V.IsDeleted {
			if _, i := b.Buf.LocateThread(stw.K); i == -1 {
				/* silly golang */
				stw.V.IsDeleted = true
				b.Intr.Active.Threads.Map(stw.K, stw.V)

				for pw := range b.Intr.Active.Posts.Range() {
					if pw.V.Parent == stw.V.No {
						pw.V.IsDeleted = true
						b.Intr.Active.Posts.Map(pw.K, pw.V)
					}
				}

				if stw.V.Page < 10 {
					c.Printf("%s/%d was deleted!\n", b.Name, stw.V.No)
				} else {
					c.Printf("%s/%d was pruned.\n", b.Name, stw.V.No)
				}
			}
		}
	}
}

func threads(c *conn.Con, b *conn.Bored) error {
	u := conn.Url(conn.Static, b.Name+"/threads.json")
	if err := c.Gather(u, &b.Buf.Threads); err != nil {
		return errs.Wrapf(err, "Failed to gather %s", u.String())
	}
	for _, v := range b.Buf.Threads {
		for _, vv := range v.Threads {
			b.Buf.TotalThreads++
			vv.LastReplies = b.Intr.Active.Threads.Find(vv.No).Replies
			vv.Page = v.Page
			b.Intr.Active.Threads.Map(vv.No, vv)

		}
	}

	/* locate deleted threads */
	deletedThreads(c, b)

	b.Buf.Threads = b.Buf.Threads[:0] /* truncate */

	return nil
}

func deletedPosts(c *conn.Con, b *conn.Bored, st conn.SmallThread) {
	for pw := range b.Intr.Active.Posts.Range() {
		if pw.V.Parent == st.No {
			if !st.IsDeleted {
				if !pw.V.IsDeleted {
					if b.Buf.Locate(pw.V.No) == -1 {
						/* silly golang */
						pw.V.IsDeleted = true
						b.Intr.Active.Posts.Map(pw.V.No, pw.V)
						c.Printf("parent: %d: post %d was deleted from %s/%d\n-- -- --\n%+v\n-- -- --\n",
							st.No, pw.V.No, b.Name, pw.V.Parent, b.Intr.Active.Posts.Find(pw.V.No))
					}
				}
			}
		}
	}
}

func posts(c *conn.Con, b *conn.Bored, st conn.SmallThread, imgChan chan conn.Image) error {
	//fmt.Println("entered posts")
	u := conn.Url(conn.Static, conn.MakePath(b.Name, st.No))
	if err := c.Gather(u, &b.Buf); err != nil {
		return errs.Wrapf(err, "Failed to gather %s", u.String())
	}

	//fmt.Println("ranging over posts")
	for _, v := range b.Buf.Posts {
		v.Parent = st.No
		b.Intr.Active.Posts.Map(v.No, v)

		/* images - not finished */
		// fmt.Println("Doing image stuff")
		if v.Tim > 0 {
			u := conn.MakeImgURL(conn.MakeImgPath(v.Ext, b.Name, conn.MakeImgN(v.Tim)))
			p := b.Intr.Active.Posts.Find(v.No)
			p.Image.URL = u
			p.Image.Original = p.Filename + p.Ext
			p.Image.Filename = conn.MakeImgN(p.Tim) + p.Ext
			b.Intr.Active.Posts.Map(v.No, p)
			// fmt.Println("NOT putting it down imgChan")
			// imgChan <- p.Image /* download */
		}
	}

	//fmt.Println("deleted posts time")
	/* locate deleted posts */
	deletedPosts(c, b, st)

	//fmt.Println("truncating b.Buf.Posts")
	b.Buf.Posts = b.Buf.Posts[:0] /* truncate */

	//fmt.Println("returning nil")
	return nil
}

/* this program has turned out to be very loopy, so this keeps
 * the number of posts/threads to be looped through to a minimum. */
func clean(c *conn.Con, b *conn.Bored) {
	//fmt.Println("starting clean")
	newActiveThreads := conn.NewEmptyIntST()
	for stw := range b.Intr.Active.Threads.Range() {
		if stw.V.IsDeleted {
			b.Intr.Deleted.Threads.Map(stw.V.No, stw.V)
		} else {
			newActiveThreads.Map(stw.V.No, stw.V)
		}
	}
	//fmt.Println("pt1 clean")
	b.Intr.Active.Threads = newActiveThreads
	newActivePosts := conn.NewEmptyIntPost()
	for pw := range b.Intr.Active.Posts.Range() {
		if pw.V.IsDeleted {
			b.Intr.Deleted.Posts.Map(pw.V.No, pw.V)
		} else {
			newActivePosts.Map(pw.V.No, pw.V)
		}
	}
	//fmt.Println("pt2 clean")
	b.Intr.Active.Posts = newActivePosts
	//c.Println("Cleaned!")
}

func sweep(c *conn.Con, b *conn.Bored, imgChan chan conn.Image) error {
	if err := threads(c, b); err != nil {
		return err
	}
	for stw := range b.Intr.Active.Threads.Range() {
		if stw.V.Replies > stw.V.LastReplies {
			c.Printf("%d new posts on %s/%d\n", stw.V.Replies-stw.V.LastReplies, b.Name, stw.V.No) /* */
			//fmt.Println("getting new posts")
			if err := posts(c, b, stw.V, imgChan); err != nil {
				c.Printf("%s\n", err)
			}
			//fmt.Println("finished getting new posts")
		}
	}
	//fmt.Println("cleaning")
	clean(c, b)
	return nil
}
