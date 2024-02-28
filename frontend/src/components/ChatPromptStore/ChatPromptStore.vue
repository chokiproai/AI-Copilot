<script setup lang="ts">
import { ref } from 'vue';
import { NModal, NList, NListItem, NButton, useMessage, NSpace, NInput, NUpload, type UploadFileInfo, NEmpty } from 'naive-ui';
import { usePromptStore, type IPrompt, type IPromptDownloadConfig } from '@/stores/modules/prompt';
import { storeToRefs } from 'pinia';
import VirtualList from 'vue3-virtual-scroll-list';
import ChatPromptItem from './ChatPromptItem.vue';

const messgae = useMessage();
const promptStore = usePromptStore();
const { promptDownloadConfig, isShowPromptSotre, promptList, keyword, searchPromptList, optPromptConfig } = storeToRefs(promptStore);

const isShowDownloadPop = ref(false);

const isImporting = ref(false);
const isExporting = ref(false);

const showAddPromptPop = () => {
  optPromptConfig.value.isShow = true;
  optPromptConfig.value.type = 'add';
  optPromptConfig.value.title = 'Thêm lời nhắc';
  optPromptConfig.value.newPrompt = {
    act: '',
    prompt: '',
  };
};

const savePrompt = () => {
  const { type, tmpPrompt, newPrompt } = optPromptConfig.value;
  if (!newPrompt.act) {
    return messgae.error('Tiêu đề từ nhắc nhở không được để trống');
  }
  if (!newPrompt.prompt) {
    return messgae.error('Mô tả từ nhắc nhở không được để trống');
  }
  if (type === 'add') {
    promptList.value = [newPrompt, ...promptList.value];
    messgae.success('Thêm lời nhắc thành công');
  } else if (type === 'edit') {
    if (newPrompt.act === tmpPrompt?.act && newPrompt.prompt === tmpPrompt?.prompt) {
      messgae.warning('Từ nhắc nhở không thay đổi');
      optPromptConfig.value.isShow = false;
      return;
    }
    const rawIndex = promptList.value.findIndex((x) => x.act === tmpPrompt?.act && x.prompt === tmpPrompt?.prompt);
    if (rawIndex > -1) {
      promptList.value[rawIndex] = newPrompt;
      messgae.success('Chỉnh sửa lời nhắc thành công');
    } else {
      messgae.error('Chỉnh sửa lỗi lời nhắc');
    }
  }
  optPromptConfig.value.isShow = false;
};

const readFile = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = function (ev) {
      resolve(ev.target?.result as string);
    };
    reader.onerror = reject;
    reader.readAsText(file);
  });
};

const importPrompt = async (options: { file: UploadFileInfo; fileList: Array<UploadFileInfo>; event?: Event }) => {
  // console.log(options.file);
  if (options.file.file) {
    isImporting.value = true;
    const fileText = await readFile(options.file.file);
    const promptData = JSON.parse(fileText);
    const result = promptStore.addPrompt(promptData);
    if (result.result) {
      messgae.info(`Tệp được tải lên có chứa ${promptData.length} Dữ liệu bài viết`);
      messgae.success(`Đã nhập thành công ${result.data?.successCount} Dữ liệu hợp lệ`);
    } else {
      messgae.error(result.msg || 'Định dạng từ nhắc nhở là sai');
    }
    isImporting.value = false;
  } else {
    messgae.error('Đã tải lên tệp sai');
  }
};

const exportPrompt = () => {
  if (promptList.value.length === 0) {
    return messgae.error('Hiện tại không có dữ liệu từ nhắc được xuất.');
  }
  isExporting.value = true;
  const jsonDataStr = JSON.stringify(promptList.value);
  const blob = new Blob([jsonDataStr], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = 'BingAIPrompts.json';
  link.click();
  URL.revokeObjectURL(url);
  messgae.success('Đã xuất thành công từ vựng lời nhắc');
  isExporting.value = false;
};

const clearPrompt = () => {
  promptList.value = [];
  messgae.success('Đã xóa từ vựng nhắc nhở thành công');
};

const downloadPrompt = async (config: IPromptDownloadConfig) => {
  if (!config.url) {
    return messgae.error('Vui lòng nhập link tải trước');
  }
  config.isDownloading = true;
  let jsonData: Array<IPrompt>;
  if (config.url.endsWith('.json')) {
    jsonData = await fetch(config.url).then((res) => res.json());
  } else if (config.url.endsWith('.csv')) {
    const csvData = await fetch(config.url).then((res) => res.text());
    console.log(csvData);
    jsonData = csvData
      .split('\n')
      .filter((x) => x)
      .map((x) => {
        const arr = x.split('","');
        return {
          act: arr[0].slice(1),
          prompt: arr[1]?.slice(1),
        };
      });
    jsonData.shift();
  } else {
    config.isDownloading = false;
    return messgae.error('Việc tải xuống các từ nhắc nhở cho hậu tố này hiện không được hỗ trợ.');
  }
  config.isDownloading = false;
  const result = promptStore.addPrompt(jsonData);
  if (result.result) {
    messgae.info(`Tệp tải xuống chứa ${jsonData.length} Dữ liệu bài viết`);
    messgae.success(`Đã nhập thành công ${result.data?.successCount} Dữ liệu hợp lệ`);
  } else {
    messgae.error(result.msg || 'Định dạng từ nhắc nhở là sai');
  }
};
</script>

<template>
  <NModal class="w-11/12 xl:w-[900px]" v-model:show="isShowPromptSotre" preset="card" title="Nhắc từ điển đồng nghĩa">
    <div class="flex justify-start flex-wrap gap-2 px-5 pb-2">
      <NInput class="basis-full xl:basis-0 xl:min-w-[300px]" placeholder="Tìm kiếm từ gợi ý" v-model:value="keyword" :clearable="true"></NInput>
      <NButton secondary type="info" @click="isShowDownloadPop = true">Tải xuống</NButton>
      <NButton secondary type="info" @click="showAddPromptPop">Thêm vào</NButton>
      <NUpload class="w-[56px] xl:w-auto" accept=".json" :default-upload="false" :show-file-list="false" @change="importPrompt">
        <NButton secondary type="success" :loading="isImporting">Nhập khẩu</NButton>
      </NUpload>
      <!-- <NButton secondary type="success">Nhập khẩu</NButton> -->
      <NButton secondary type="success" @click="exportPrompt" :loading="isExporting">Xuất khẩu</NButton>
      <NButton secondary type="error" @click="clearPrompt">Thông thoáng</NButton>
    </div>
    <VirtualList
      v-if="searchPromptList.length > 0"
      class="h-[40vh] xl:h-[60vh] overflow-y-auto"
      :data-key="'prompt'"
      :data-sources="searchPromptList"
      :data-component="ChatPromptItem"
      :keeps="10"
    />
    <NEmpty v-else class="h-[40vh] xl:h-[60vh] flex justify-center items-center" description="Không có dữ liệu">
      <template #extra>
        <NButton secondary type="info" @click="isShowDownloadPop = true">Tải xuống các từ gợi ý</NButton>
      </template>
    </NEmpty>
  </NModal>
  <NModal class="w-11/12 xl:w-[600px]" v-model:show="optPromptConfig.isShow" preset="card" :title="optPromptConfig.title">
    <NSpace vertical>
      tiêu đề
      <NInput placeholder="Vui lòng nhập tiêu đề" v-model:value="optPromptConfig.newPrompt.act"></NInput>
      mô tả
      <NInput placeholder="Vui lòng nhập mô tả" type="textarea" v-model:value="optPromptConfig.newPrompt.prompt"></NInput>
      <NButton block secondary type="info" @click="savePrompt">Lưu</NButton>
    </NSpace>
  </NModal>
  <NModal class="w-11/12 xl:w-[600px]" v-model:show="isShowDownloadPop" preset="card" title="Tải xuống các từ gợi ý">
    <NList class="overflow-y-auto rounded-lg" hoverable clickable>
      <NListItem v-for="(config, index) in promptDownloadConfig" :key="index">
        <a v-if="config.type === 1" class="no-underline text-blue-500" :href="config.url" target="_blank" rel="noopener noreferrer">{{ config.name }}</a>
        <NInput v-else-if="config.type === 2" placeholder="Vui lòng nhập link tải, hỗ trợ json và csv " v-model:value="config.url"></NInput>
        <template #suffix>
          <div class="flex justify-center gap-5">
            <a class="no-underline" v-if="config.type === 1" :href="config.refer" target="_blank" rel="noopener noreferrer">
              <NButton secondary>nguồn</NButton>
            </a>
            <NButton secondary type="info" @click="downloadPrompt(config)" :loading="config.isDownloading">Tải xuống</NButton>
          </div>
        </template>
      </NListItem>
    </NList>
  </NModal>
</template>
