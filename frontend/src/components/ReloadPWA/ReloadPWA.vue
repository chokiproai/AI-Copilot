<script setup lang="ts">
import { ref } from 'vue';
import { NModal, NButton, useMessage } from 'naive-ui';
import { useRegisterSW } from 'virtual:pwa-register/vue';

const message = useMessage();
const isShowModal = ref(false);
const isUpdateLoading = ref(false);

const { offlineReady, needRefresh, updateServiceWorker } = useRegisterSW({
  immediate: true,
  onRegisteredSW(swUrl, r) {
    // console.log(`Service Worker at: ${swUrl}`);
    r &&
      setInterval(async () => {
        // console.log('Checking for sw update');
        await r.update();
      }, 60 * 1000);
  },
  onNeedRefresh() {
    isShowModal.value = true;
    // updateServiceWorker();
  },
});

const close = async () => {
  offlineReady.value = false;
  needRefresh.value = false;
  isShowModal.value = false;
};

const update = async () => {
  isUpdateLoading.value = true;
  await updateServiceWorker();
  message.success('Đã chuyển sang phiên bản mới');
  isUpdateLoading.value = false;
  isShowModal.value = false;
};
</script>

<template>
  <NModal v-model:show="isShowModal" preset="dialog" title="Cập nhật mẹo" content="Phiên bản mới đã được cập nhật. Bạn có muốn chuyển đổi không?">
    <template #action>
      <NButton size="large" @click="close">Hẹn gặp lại lần sau</NButton>
      <NButton ghost size="large" type="info" @click="update" :loading="isUpdateLoading">Chuyển đổi ngay bây giờ</NButton>
    </template>
  </NModal>
</template>
