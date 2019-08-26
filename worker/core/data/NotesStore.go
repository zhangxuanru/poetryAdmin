package data

import (
	"poetryAdmin/worker/app/models"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/define"
	"time"
)

//保存详情文本信息
type NotesStore struct {
}

func NewNotesStore() *NotesStore {
	return new(NotesStore)
}

//保存详情文本信息
func (a *NotesStore) SaveNotes(content *define.ContentData) (id int64, err error) {
	content.Content = tools.TrimDivHtml(content.Content)
	notes := &models.Notes{
		Title:      content.Title,
		Content:    content.Content,
		PlayUrl:    content.PlayUrl,
		PlaySrcUrl: content.PlaySrcUrl,
		HtmlSrcUrl: content.HtmlSrcUrl,
		Type:       content.Type,
		Introd:     content.Introd,
		AddDate:    time.Now().Unix(),
		UpdateDate: time.Now().Unix(),
	}
	if content.Id > 0 {
		notes.Id = content.Id
	}
	if len(content.FileName) > 0 {
		notes.FileName = content.FileName
	}
	id, err = models.NewNotes().SaveNotes(notes)
	return
}

//更新MP3文件路径
func (a *NotesStore) UpdateMp3Path(content *define.ContentData) (id int64, err error) {
	store := &models.Notes{
		Id:       content.Id,
		FileName: content.FileName,
	}
	id, err = models.NewNotes().UpdateNotes(store, "file_name")
	return
}
