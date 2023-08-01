#if defined(_WIN32) || defined(_WIN64)

#include "filepicker.h"

#include <string.h>
#include <string>

#include <Windows.h>
#include <ShObjIdl.h>

static bool file_dialog(const FILEOPENDIALOGOPTIONS opt, const std::string &title, std::string &path, bool save_dialog = false)
{
    // Initialize this piece of lib
    HRESULT hr = CoInitializeEx(NULL, COINIT_APARTMENTTHREADED | COINIT_DISABLE_OLE1DDE);
    if (!SUCCEEDED(hr))
        return false;
    defer([]
          { CoUninitialize(); });

    // Create file chooser dialog
    IFileDialog *file_dialog;
    if (save_dialog)
    {
        hr = CoCreateInstance(CLSID_FileSaveDialog, NULL, CLSCTX_ALL, IID_IFileSaveDialog, (void **)&file_dialog);
    }
    else
    {
        hr = CoCreateInstance(CLSID_FileOpenDialog, NULL, CLSCTX_ALL, IID_IFileOpenDialog, (void **)&file_dialog);
    }
    if (!SUCCEEDED(hr))
        return false;
    defer([file_dialog]
          { file_dialog->Release(); });

    // Show dialog and get file info
    std::wstring wtitle(title.begin(), title.end());
    file_dialog->SetOptions(opt);
    file_dialog->SetTitle(wtitle.c_str());
    if (save_dialog && path.length() > 0)
    {
        std::wstring wpath(path.begin(), path.end());
        file_dialog->SetFileName(wpath.c_str());

        auto ch_index = wpath.find_last_of(L'.');
        ch_index = ch_index == std::string::npos ? 0 : ch_index;
        std::wstring filetype = wpath.substr(ch_index, wpath.size() - ch_index);
        filetype = L'*' + filetype;

        COMDLG_FILTERSPEC filetypes_filter[] = {
            {L"Current type", filetype.c_str()},
            {L"Any", L"*.*"},
        };
        file_dialog->SetFileTypes(2, filetypes_filter);
    }
    hr = file_dialog->Show(nullptr);
    if (!SUCCEEDED(hr))
        return false;

    IShellItem *selected_item;
    hr = file_dialog->GetResult(&selected_item);
    if (!SUCCEEDED(hr))
        return false;
    defer([selected_item]
          { selected_item->Release(); });

    // Get file path
    wchar_t *path_ptr;
    hr = selected_item->GetDisplayName(SIGDN_FILESYSPATH, &path_ptr);
    if (!SUCCEEDED(hr))
        return false;
    defer([path_ptr]
          { CoTaskMemFree(path_ptr); });

    std::wstring wpath(path_ptr);
    path = std::string(wpath.begin(), wpath.end());

    return true;
}

bool open_directory(const char *title, char *out_path)
{
    std::string path;
    bool suc = file_dialog(FOS_PICKFOLDERS | FOS_PATHMUSTEXIST | FOS_FILEMUSTEXIST, title, path, false);
    if (suc)
    {
        strcpy(out_path, path.c_str());
    }
    return suc;
}

bool open_file(const char *title, char *out_path)
{
    std::string path;
    bool suc = file_dialog(FOS_PATHMUSTEXIST | FOS_FILEMUSTEXIST, title, path, false);
    if (suc)
    {
        strcpy(out_path, path.c_str());
    }
    return suc;
}

bool save_file(const char *title, char *out_path)
{
    std::string path(out_path);
    bool suc = file_dialog(FOS_OVERWRITEPROMPT | FOS_PATHMUSTEXIST, title, path, true);
    if (suc)
    {
        strcpy(out_path, path.c_str());
    }
    return suc;
}

#endif