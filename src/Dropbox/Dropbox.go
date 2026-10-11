package Dropbox

import (
	"Helpers"
	"Logger"
	"fmt"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/dustin/go-humanize"
	"github.com/mitchellh/ioprogress"
	"io"
	"os"
	"path"
	"time"
)

var Logging Logger.Interface

const chunkSize int64 = 1 << 24

func uploadChunked(dbx files.Client, r io.Reader, commitInfo *files.CommitInfo, sizeTotal int64) (err error) {
	res, err := dbx.UploadSessionStart(files.NewUploadSessionStartArg(),
		&io.LimitedReader{R: r, N: chunkSize})
	if err != nil {
		return err
	}

	written := chunkSize

	for (sizeTotal - written) > chunkSize {
		cursor := files.NewUploadSessionCursor(res.SessionId, uint64(written))
		args := files.NewUploadSessionAppendArg(cursor)

		err = dbx.UploadSessionAppendV2(args, &io.LimitedReader{R: r, N: chunkSize})
		if err != nil {
			return err
		}
		written += chunkSize
	}

	cursor := files.NewUploadSessionCursor(res.SessionId, uint64(written))
	args := files.NewUploadSessionFinishArg(cursor, commitInfo)

	if _, err = dbx.UploadSessionFinish(args, r); err != nil {
		return err
	}

	return err
}

func UploadFile(token string, src string, uploadPath string) bool {
	config := dropbox.Config{
		Token:    token,
		LogLevel: dropbox.LogOff, // if needed, set the desired logging level. Default is off
	}
	dst := "/" + uploadPath + "/" + path.Base(src)
	dst = Helpers.ReplaceAll(dst, "//", "/")

	contents, err := os.Open(src)
	if err != nil {
		return false
	}
	defer func() {
		_ = contents.Close()
	}()

	contentsInfo, err := contents.Stat()
	if err != nil {
		return false
	}

	commitInfo := files.NewCommitInfo(dst)
	commitInfo.Mode.Tag = "overwrite"

	// The Dropbox API only accepts timestamps in UTC with second precision.
	commitInfo.ClientModified = time.Now().UTC().Round(time.Second)
	progressbar := &ioprogress.Reader{
		Reader: contents,
		DrawFunc: ioprogress.DrawTerminalf(os.Stderr, func(progress, total int64) string {
			return fmt.Sprintf("%s/%s",
				humanize.IBytes(uint64(progress)), humanize.IBytes(uint64(total)))
		}),
		Size: contentsInfo.Size(),
	}

	dbx := files.New(config)
	if contentsInfo.Size() > chunkSize {
		if err = uploadChunked(dbx, progressbar, commitInfo, contentsInfo.Size()); err != nil {
			return false
		}
	} else {
		if _, err = dbx.Upload(commitInfo, progressbar); err != nil {
			return false
		}
	}
	return true
}
