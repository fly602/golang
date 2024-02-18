#include <stdio.h>  
#include <apt-pkg/pkg.h>  
#include <apt-pkg/error.h>  
  
int main(int argc, char *argv[]) {  
    // 初始化库  
    pkgInitConfig(NULL);  
    pkgLoadConfig(NULL);  
      
    // 创建一个新的包状态对象  
    pkgCache::State *state = pkgCache::State::Create();  
    if (!state) {  
        fprintf(stderr, "无法创建包状态对象\n");  
        return 1;  
    }  
      
    // 更新软件包列表  
    if (!state->Update(false)) {  
        fprintf(stderr, "更新软件包列表失败\n");  
        return 1;  
    }  
      
    // 遍历所有软件包并检查更新  
    for (auto pkg = state->PkgBegin(); pkg != state->PkgEnd(); ++pkg) {  
        if (pkg->Upgradable()) {  
            printf("软件包 '%s' 有一个可用的更新\n", pkg->Name());  
        } else {  
            printf("软件包 '%s' 已经是最新版本\n", pkg->Name());  
        }  
    }  
      
    // 释放资源  
    delete state;  
      
    return 0;  
}