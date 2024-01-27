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
  optPromptConfig.value.title = 'Thêm từ gợi ý';
  optPromptConfig.value.newPrompt = {
    act: '',
    prompt: '',
  };
};

const savePrompt = () => {
  const { type, tmpPrompt, newPrompt } = optPromptConfig.value;
  if (!newPrompt.act) {
    return messgae.error('Tiêu đề từ gợi ý không được để trống');
  }
  if (!newPrompt.prompt) {
    return messgae.error('Mô tả từ gợi ý không được để trống');
  }
  if (type === 'add') {
    promptList.value = [newPrompt, ...promptList.value];
    messgae.success('Đã thêm từ gợi ý thành công');
  } else if (type === 'edit') {
    if (newPrompt.act === tmpPrompt?.act && newPrompt.prompt === tmpPrompt?.prompt) {
      messgae.warning('Từ gợi ý không thay đổi');
      optPromptConfig.value.isShow = false;
      return;
    }
    const rawIndex = promptList.value.findIndex((x) => x.act === tmpPrompt?.act && x.prompt === tmpPrompt?.prompt);
    if (rawIndex > -1) {
      promptList.value[rawIndex] = newPrompt;
      messgae.success('Đã chỉnh sửa từ gợi ý thành công');
    } else {
      messgae.error('Chỉnh sửa từ gợi ý bị lỗi');
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
      messgae.info(`Tệp tải lên chứa ${promptData.length} dữ liệu`);
      messgae.success(`Đã nhập thành công ${result.data?.successCount} dữ liệu hợp lệ`);
    } else {
      messgae.error(result.msg || 'Định dạng từ gợi ý có lỗi');
    }
    isImporting.value = false;
  } else {
    messgae.error('Tệp tải lên có lỗi');
  }
};

const exportPrompt = () => {
  if (promptList.value.length === 0) {
    return messgae.error('Không có dữ liệu từ gợi ý để xuất');
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
  messgae.success('Đã xuất kho từ gợi ý thành công');
  isExporting.value = false;
};

const clearPrompt = () => {
  promptList.value = [];
  messgae.success('Đã xóa sạch kho từ gợi ý');
};

const downloadPrompt = async (config: IPromptDownloadConfig) => {
  if (!config.url) {
    return messgae.error('Vui lòng nhập liên kết tải xuống trước');
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
  return messgae.error('Tạm thời không hỗ trợ tải xuống từ gợi ý với phần mở rộng này');
  }
  config.isDownloading = false;
  const result = promptStore.addPrompt(jsonData);
  if (result.result) {
    messgae.info(`Tệp tải xuống chứa ${jsonData.length} dữ liệu`);
    messgae.success(`Đã nhập thành công ${result.data?.successCount} dữ liệu hợp lệ`);
  } else {
    messgae.error(result.msg || 'Định dạng từ gợi ý có lỗi');
  }
  };
  </script>

<template>
  <NModal class="w-11/12 xl:w-[900px]" v-model:show="isShowPromptSotre" preset="card" title="Kho từ gợi ý">
    <div class="flex justify-start flex-wrap gap-2 px-5 pb-2">
      <NInput class="basis-full xl:basis-0 xl:min-w-[300px]" placeholder="Tìm kiếm từ gợi ý" v-model:value="keyword" :clearable="true"></NInput>
      <NButton secondary type="info" @click="isShowDownloadPop = true">Tải xuống</NButton>
      <NButton secondary type="info" @click="showAddPromptPop">Thêm</NButton>
      <NUpload class="w-[56px] xl:w-auto" accept=".json" :default-upload="false" :show-file-list="false" @change="importPrompt">
        <NButton secondary type="success" :loading="isImporting">Nhập</NButton>
      </NUpload>
      <!-- <NButton secondary type="success">Nhập</NButton> -->
      <NButton secondary type="success" @click="exportPrompt" :loading="isExporting">Xuất</NButton>
      <NButton secondary type="error" @click="clearPrompt">Xóa sạch</NButton>
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
        <NButton secondary type="info" @click="isShowDownloadPop = true">Tải xuống từ gợi ý</NButton>
      </template>
    </NEmpty>
  </NModal>
  <NModal class="w-11/12 xl:w-[600px]" v-model:show="optPromptConfig.isShow" preset="card" :title="optPromptConfig.title">
    <NSpace vertical>
      Tiêu đề
      <NInput placeholder="Nhập tiêu đề" v-model:value="optPromptConfig.newPrompt.act"></NInput>
      Mô tả
      <NInput placeholder="Nhập mô tả" type="textarea" v-model:value="optPromptConfig.newPrompt.prompt"></NInput>
      <NButton block secondary type="info" @click="savePrompt">Lưu</NButton>
    </NSpace>
  </NModal>
  <NModal class="w-11/12 xl:w-[600px]" v-model:show="isShowDownloadPop" preset="card" title="Tải xuống từ gợi ý">
    <NList class="overflow-y-auto rounded-lg" hoverable clickable>
      <NListItem v-for="(config, index) in promptDownloadConfig" :key="index">
        <a v-if="config.type === 1" class="no-underline text-blue-500" :href="config.url" target="_blank" rel="noopener noreferrer">{{ config.name }}</a>
        <NInput v-else-if="config.type === 2" placeholder="Nhập liên kết tải xuống, hỗ trợ json và csv " v-model:value="config.url"></NInput>
        <template #suffix>
          <div class="flex justify-center gap-5">
            <a class="no-underline" v-if="config.type === 1" :href="config.refer" target="_blank" rel="noopener noreferrer">
              <NButton secondary>Nguồn</NButton>
            </a>
            <NButton secondary type="info" @click="downloadPrompt(config)" :loading="config.isDownloading">Tải xuống</NButton>
          </div>
        </template>
      </NListItem>
    </NList>
  </NModal>
</template>