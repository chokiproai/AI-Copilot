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
    return message.error('Tiêu đề từ gợi ý không được để trống');
  }
  if (!newPrompt.prompt) {
    return message.error('Mô tả từ gợi ý không được để trống');
  }
  if (type === 'add') {
    promptList.value = [newPrompt, ...promptList.value];
    message.success('Thêm từ gợi ý thành công');
  } else if (type === 'edit') {
    if (newPrompt.act === tmpPrompt?.act && newPrompt.prompt === tmpPrompt?.prompt) {
      message.warning('Từ gợi ý không thay đổi');
      optPromptConfig.value.isShow = false;
      return;
    }
    const rawIndex = promptList.value.findIndex((x) => x.act === tmpPrompt?.act && x.prompt === tmpPrompt?.prompt);
    if (rawIndex > -1) {
      promptList.value[rawIndex] = newPrompt;
      message.success('Chỉnh sửa từ gợi ý thành công');
    } else {
      message.error('Lỗi khi chỉnh sửa từ gợi ý');
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

const importPrompt = async (options) => {
  // console.log(options.file);
  if (options.file.file) {
    isImporting.value = true;
    const fileText = await readFile(options.file.file);
    const promptData = JSON.parse(fileText);
    const result = promptStore.addPrompt(promptData);
    if (result.result) {
      message.info(`Tập tin tải lên chứa ${promptData.length} dữ liệu`);
      message.success(`Nhập thành công ${result.data?.successCount} dữ liệu hợp lệ`);
    } else {
      message.error(result.msg || 'Định dạng từ gợi ý không chính xác');
    }
    isImporting.value = false;
  } else {
    message.error('Tập tin tải lên không hợp lệ');
  }
};

const exportPrompt = () => {
  if (promptList.value.length === 0) {
    return message.error('Hiện không có dữ liệu từ gợi ý để xuất');
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
  message.success('Xuất bộ từ gợi ý thành công');
  isExporting.value = false;
};

const clearPrompt = () => {
  promptList.value = [];
  message.success('Đã xóa bộ từ gợi ý thành công');
};

const downloadPrompt = async (config) => {
  if (!config.url) {
    return message.error('Vui lòng nhập đường dẫn tải xuống trước');
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
    return message.error('Hiện không hỗ trợ tải xuống các tệp có phần mở rộng này');
  }
  config.isDownloading = false;
  const result = promptStore.addPrompt(jsonData);
  if (result.result) {
    message.info(`Tập tin tải xuống chứa ${jsonData.length} dữ liệu`);
    message.success(`Nhập thành công ${result.data?.successCount} dữ liệu hợp lệ`);
  } else {
    message.error(result.msg || 'Định dạng từ gợi ý không chính xác');
  }
};
</script>

<template>
  <NModal class="w-11/12 xl:w-[900px]" v-model:show="isShowPromptSotre" preset="card" title="Từ gợi ý">
    <div class="flex justify-start flex-wrap gap-2 px-5 pb-2">
      <NInput class="basis-full xl:basis-0 xl:min-w-[300px]" placeholder="Tìm kiếm từ gợi ý" v-model:value="keyword" :clearable="true"></NInput>
      <NButton secondary type="info" @click="isShowDownloadPop = true">Tải xuống</NButton>
      <NButton secondary type="info" @click="showAddPromptPop">Thêm</NButton>
      <NUpload class="w-[56px] xl:w-auto" accept=".json" :default-upload="false" :show-file-list="false" @change="importPrompt">
        <NButton secondary type="success" :loading="isImporting">Nhập</NButton>
      </NUpload>
      <!-- <NButton secondary type="success">Nhập</NButton> -->
      <NButton secondary type="success" @click="exportPrompt" :loading="isExporting">Xuất</NButton>
      <NButton secondary type="error" @click="clearPrompt">Xóa</NButton>
    </div>
    <VirtualList
      v-if="searchPromptList.length > 0"
      class="h-[40vh] xl:h-[60vh] overflow-y-auto"
      :data-key="'prompt'"
      :data-sources="searchPromptList"
      :data-component="ChatPromptItem"
      :keeps="10"
    />
<NEmpty v-else class="h-[40vh] xl:h-[60vh] flex justify-center items-center" description="Hiện không có dữ liệu">
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
      <NInput v-else-if="config.type === 2" placeholder="Vui lòng nhập đường dẫn tải xuống, hỗ trợ json và csv" v-model:value="config.url"></NInput>
      <template #suffix>
        <div class="flex justify-center gap-5">
          <a class="no-underline" v-if="config.type === 1" :href="config.refer" target="_blank" rel="noopener noreferrer">
            <NButton secondary>Xuất xứ</NButton>
          </a>
          <NButton secondary type="info" @click="downloadPrompt(config)" :loading="config.isDownloading">Tải xuống</NButton>
        </div>
      </template>
    </NListItem>
  </NList>
</NModal>
</template>
