package scanner

import (
	"github.com/smook1980/medialocker"
	"github.com/smook1980/medialocker/cli"
	. "github.com/smook1980/medialocker/models"
	"github.com/smook1980/medialocker/types"
	"github.com/smook1980/medialocker/util"
	cli2 "github.com/urfave/cli"
)

var Command = cli.Command{
	Name:        "scanner",
	Usage:       "Scan for media at PATH.",
	Description: "Scans path for image and video media.",
	Action: func(c *cli2.Context) error {
		app, errs := medialocker.NewAppBuilder().WithConfiguration(medialocker.FileConfiguration("")).Build()

		if len(errs) != 0 {
			return util.MultiError(errs...)
		}

		var db *medialocker.DBConnection
		var err error

		if db, err = app.Registry.DB(); err != nil {
			app.Log.Fatalf("Failed to open db! %s", err)
		}

		defer db.Close()

		scanner := NewScanner("/Users/smook/Downloads", 8, app.Log)

		scanner.Each(func(mp types.MediaPath) {
			fp, _ := NewFilePath(mp.Realpath, mp.Hash)
			tx := db.Begin()
			tx.
				Where(FilePath{Basename: fp.Basename, Dirname: fp.Dirname}).
				Attrs(fp).
				FirstOrCreate(&fp)

			tx.Commit()

			app.Log.
				WithField("prefix", "media_path").
				WithField("event", "media_path_found").
				WithField("command", "scanner").
				Debugf("%+v", fp)
		})

		scanner.Each(func(mp types.MediaPath) {
			app.Log.
				WithField("prefix", "scanner_cmd").
				WithField("event", "media_path_info").
				WithField("Hash", mp.Hash).
				WithField("MediaType", mp.Type.String()).
				WithField("Realpath", mp.Realpath).
				Infoln("Located MediaPath")
		})

		app.Start("Scanner", scanner.Module)
		app.Wait()
		app.Shutdown()

		return nil
	},
}
