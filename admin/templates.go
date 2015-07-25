package admin

import (
	"os"
	"path/filepath"
	"strings"
)

type FolderContent struct {
	File string `json:"file"`
	Type string `json:"type"`
	Path string `json:"path"`
}

func ListTemplates(w *wrapper.Wrapper) {
	rootdirectory := filepath.Join(w.SiteConfig.Directory, w.SiteConfigs.AssetsDirectory, "templates")
	var workingdirectory string
	list = make([]FolderContent, 0)
	if len(w.APIParams) > 0 {
		workingdirectory = strings.Join(w.APIParams)
		// This is not top level so add the up directory
		fc := FolderContent{"..", "directory", filepath.Join(workingdirectory, "/..")}
		list = append(list, fc)
	}
	glob := filepath.Join(rootdirectory, workingdirectory)
	abs, err := filepath.Abs(glob)
	// Check to insure this is within the site directories
	if !HasPrefix(abs, rootdirectory) || err != nil {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	glob = filepath.Join(glob, "*")
	files, err := filepath.Glob(glob)
	for _, file := range files {
		info, _ := os.Stat(file)
		if info.IsDir() {
			fc := FolderContent{filepath.Base(file), "template", filepath.Join(workingdirectory, filepath.Base(file))}
			list = append(list, fc)
		} else {
			if filepath.Ext(file) == "html" {
				fc := FolderContent{filepath.Base(file), "template", filepath.Join(workingdirectory, filepath.Base(file))}
				list = append(list, fc)
			}
		}
	}
	w.SetTemplate("admin/template_list.html")
	w.SetPayload("folder_content", list)
	w.SetPayload("working_directory", workingdirectory)
	w.Serve()
	return
}

func TemplateEditor(w *wrapper.Wrapper) {
	if len(w.APIParams) == 0 {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	if w.Request.Method != "POST" {
		ContentTypeEditorForm(w)
		return
	}
	ContentTypeEditorSubmit(w)
	return

}

func TemplateEditorForm(w *wrapper.Wrapper) {

}

func TemplateEditorSubmit(w *wrapper.Wrapper) {

}
