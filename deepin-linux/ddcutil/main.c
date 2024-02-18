#include <ddcutil_c_api.h>
#include <stddef.h>

int main() {
    // 获取显示器信息列表
    DDCA_Display_Info_List *displayInfoList = ddca_get_display_info_list();
    if (displayInfoList == NULL) {
        printf("===>>>NULL\n");
        return 0;
    }
    for (int i = 0;i<displayInfoList->ct;i++){
        printf("===>>>model_name=%s\n",displayInfoList->info[i].model_name);
        printf("===>>>sn=%s\n",displayInfoList->info[i].sn);
        printf("===>>>sn=%s\n",(char *)displayInfoList->info[i].edid_bytes);
    }

    // 释放显示器信息列表的内存
    ddca_free_display_info_list(displayInfoList);

    return 0;
}
