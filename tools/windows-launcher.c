typedef unsigned long DWORD;
typedef int BOOL;
typedef unsigned int UINT;
__declspec(dllimport) DWORD __stdcall GetModuleFileNameA(void*, char*, DWORD);
__declspec(dllimport) BOOL __stdcall SetCurrentDirectoryA(const char*);
__declspec(dllimport) UINT __stdcall WinExec(const char*, UINT);
__declspec(dllimport) void __stdcall ExitProcess(UINT);

void mainCRTStartup(void) {
    char path[512];
    DWORD n = GetModuleFileNameA(0, path, sizeof(path));
    if (n > 0 && n < sizeof(path)) {
        for (DWORD i = n; i > 0; --i) {
            if (path[i-1] == '\\' || path[i-1] == '/') {
                path[i-1] = 0;
                SetCurrentDirectoryA(path);
                break;
            }
        }
    }
    UINT rc = WinExec("gzdoom.exe -iwad DOOM2.WAD -file bahon-shooter.pk3 +map MAP01", 1);
    ExitProcess(rc > 31 ? 0 : rc);
}
