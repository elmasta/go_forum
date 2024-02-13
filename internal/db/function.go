package forum

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	middle "forum/pkg/middleware"
	models "forum/pkg/models"
)

type Img struct {
	file       multipart.File
	fileheader *multipart.FileHeader
}

func InitImg(r *http.Request) (Img, error) {
	var img Img
	file, fileheader, err := r.FormFile("myFile")
	if err != nil {
		if err.Error() == "http: no such file" {
			return img, err
		}
		fmt.Println("formfile :", err.Error())
		return img, err
	}
	img.file = file
	img.fileheader = fileheader
	tab := strings.Split(img.fileheader.Filename, ".")
	// fileheader.Size in octer
	if len(tab) != 2 || img.fileheader.Size > (20*1000000) {
		// file.Close()
		// fmt.Println("badformat or >20Mo")
		return img, fmt.Errorf("badformat or > 20Mo")
	}
	key := tab[1]
	switch key {
	case "svg", "gif", "png", "jpeg", "jpg":
	default:
		// file.Close()
		// fmt.Println("badformat")
		return img, fmt.Errorf("badformat or > 20Mo")
	}
	return img, nil
}

func mkdirPostAsset(app *App_db, idpost int64, post *models.Post, file Img) {
	dirfile := "web/static/upload/img/post" + strconv.Itoa(int(idpost)) + "/"
	if _, err := os.Stat(dirfile); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(dirfile+"comment", os.ModePerm)
		if err != nil {
			fmt.Println("mkfirPostAsset : ", err)
			return
		}
	}
	// file, fileheader, err := r.FormFile("myFile")
	// if err != nil {
	// 	if err.Error() == "http: no such file" {
	// 		return
	// 	}
	// 	fmt.Println("formfile :", err.Error())
	// 	return
	// }
	// defer file.Close()
	//.svg,.gif,.pnj,.jpeg
	// tab := strings.Split(file.fileheader.Filename, ".")
	// fileheader.Size in octer
	// if len(tab) != 2 || file.fileheader.Size > (20*1000000) {
	// 	// file.Close()
	// 	fmt.Println("badformat or >20Mo")
	// 	return
	// }
	// key := tab[1]
	// switch key {
	// case "svg", "gif", "pnj", "jpeg", "jpg":
	// default:
	// 	// file.Close()
	// 	fmt.Println("badformat")
	// 	return
	// }
	f, err := os.OpenFile(dirfile+file.fileheader.Filename, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	defer file.file.Close()
	// fmt.Println(f)
	io.Copy(f, file.file)
	post.Img = file.fileheader.Filename
	middle.UpdateImgPoste(app.DB, idpost, file.fileheader.Filename)
}

func linkPost(app *App_db, post_id int64) (tab_like map[int64]bool, tab_dislike map[int64]bool) {
	tab_like, tab_dislike = make(map[int64]bool), make(map[int64]bool)
	rows, err := app.DB.Query("SELECT user_id,likes FROM link_post WHERE post_id = ?", post_id)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	var user_id int64
	var like bool
	for rows.Next() {
		rows.Scan(&user_id, &like)
		if like {
			tab_like[user_id] = true
		} else {
			tab_dislike[user_id] = true
		}
	}
	return
}
