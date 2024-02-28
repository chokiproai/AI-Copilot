<script setup lang="ts">
import { ref } from 'vue';
import { NButton, NEmpty, NInput, NModal, useMessage } from 'naive-ui';
import { computed } from 'vue';

const props = defineProps<{
  show: boolean;
}>();

const emit = defineEmits<{
  'update:show': [value: boolean];
}>();

const message = useMessage();

const DEMO_KEYWORD = 'Một chú mèo con đi xe máy, phóng nhanh trên đường, cảnh hoạt hình, chi tiết. ';
const keyword = ref('');
const iframeSrc = ref('');
const isCreating = ref(false);

const isShowModal = computed({
  get: () => props.show,
  set: (value) => emit('update:show', value),
});

const createImage = () => {
  if (!keyword.value) {
    message.error('Vui lòng nhập từ khóa trước');
    return;
  }
  isCreating.value = true;
  iframeSrc.value = `/images/create?re=1&showselective=1&sude=1&kseed=7500&SFX=2&q=${encodeURIComponent(keyword.value)}&t=${Date.now()}`;
};

const onClose = () => {
  keyword.value = '';
  iframeSrc.value = '';
};

const useDemo = () => {
  keyword.value = DEMO_KEYWORD;
  return createImage();
};
</script>

<template>
  <NModal class="w-11/12 lg:w-[540px] select-none" v-model:show="isShowModal" :on-close="onClose" preset="card" title="Tạo hình ảnh">
    <head class="flex justify-center gap-3 px-8">
      <NInput class="flex-1" placeholder="Lời nhắc" v-model:value="keyword" :clearable="true" @keyup.enter="createImage" maxlength="100"></NInput>
      <NButton secondary type="info" @click="createImage" :loading="isCreating">Bắt đầu tạo</NButton>
    </head>
    <main class="flex justify-center items-center">
      <iframe v-if="iframeSrc" class="w-[310px] h-[350px] xl:w-[475px] xl:h-[520px] my-4" :src="iframeSrc" frameborder="0" @load="isCreating = false"></iframe>
      <NEmpty v-else class="h-[40vh] xl:h-[60vh] flex justify-center items-center" description="Chưa được tạo">
        <template #extra>
          <NButton secondary type="info" @click="useDemo">Tạo bằng cách sử dụng ví dụ</NButton>
          <div class="text-[#c2c2c2] px-2 xl:px-10">
            <p class="text-left">Mẹo: Tính từ + danh từ + động từ + phong cách. Mô tả càng chi tiết thì hiệu quả càng tốt. </p>
            <p class="text-left">Ví dụ：{{ DEMO_KEYWORD }}</p>
          </div>
        </template>
      </NEmpty>
    </main>
  </NModal>
</template>
