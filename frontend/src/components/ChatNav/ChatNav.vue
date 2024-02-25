<script setup lang="ts">
import { h, ref, onMounted, inject, defineComponent, render } from 'vue';
import { NDropdown, type DropdownOption, NModal, NInput, NInputNumber, NButton, NGrid, NGridItem, useMessage, NImage, NForm, NFormItem, NSwitch, NTag, NSelect, NSpin, NP, NA, NConfigProvider, lightTheme, darkTheme, useOsTheme, type GlobalTheme } from 'naive-ui';
import conversationCssText from '@/assets/css/conversation.css?raw';
import settingSvgUrl from '@/assets/img/setting.svg?url';
import { usePromptStore } from '@/stores/modules/prompt';
import { storeToRefs } from 'pinia';
import ChatNavItem from './ChatNavItem.vue';
import type { Component } from 'vue';
import { isMobile } from '@/utils/utils';
import CreateImage from '@/components/CreateImage/CreateImage.vue';
import { useChatStore } from '@/stores/modules/chat';
import { useUserStore } from '@/stores/modules/user';

const isShowMore = ref(false);
const isShowSettingModal = ref(false);
const isShowAdvancedSettingModal = ref(false);
const isShowSetAboutModal = ref(false);
const isShowCookieModal = ref(false);
const isShowLoginModal = ref(false);
const isShowIframe = ref(false);
const userToken = ref('');
const userKievRPSSecAuth = ref('');
const userMUID = ref('');
const userRwBf = ref('');
const message = useMessage();
const promptStore = usePromptStore();
const { isShowPromptSotre } = storeToRefs(promptStore);
const isShowClearCacheModal = ref(false);
const isShowCreateImageModal = ref(false);
const chatStore = useChatStore();
const { isShowChatServiceSelectModal } = storeToRefs(chatStore);
const userStore = useUserStore();
const localVersion = __APP_INFO__.version;
const lastVersion = ref('...');
const { historyEnable, themeMode, uiVersion, fullCookiesEnable, cookiesStr, enterpriseEnable, customChatNum, gpt4tEnable, sydneyEnable, sydneyPrompt, passServer } = storeToRefs(userStore);

let cookiesEnable = ref(false);
let cookies = ref('');
let history = ref(true);
let themeModeSetting = ref('auto');
let uiVersionSetting = ref('v3');
let theme = ref(inject('theme'));

let settingIconStyle = ref({
  filter: 'invert(70%)',
})
let passingCFChallenge = ref(false);
const enterpriseSetting = ref(false);
const customChatNumSetting = ref(0);
const gpt4tSetting = ref(true);
const sydneySetting = ref(false);
const sydneyPromptSetting = ref('');
const passServerSetting = ref('');
const getCookieTip = ref('Đang lấy Cookie, vui lòng chờ...');

const GetLastVersion = async () => {
  const res = await fetch('https://api.github.com/repos/chokiproai/AI-Copilot/releases/latest');
  const json = await res.json();
  lastVersion.value = json.tag_name;
};

const navType = {
  login: 'login',
  setting: 'setting',
  chat: 'chat',
  notebook: 'notebook',
  compose: 'compose',
  createImage: 'createImage',
  reset: 'reset',
  about: 'about',
};
let navConfigs = ref([
  {
    key: navType.setting,
    label: 'Cài đặt',
  },
  {
    key: navType.notebook,
    label: 'Sổ tay',
  },
  {
    key: navType.reset,
    label: 'Đặt lại một lần',
  },
  {
    key: navType.about,
    label: 'Giới thiệu'
  },
]);

const themeModeOptions = ref([
  {
    label: 'Sáng',
    value: 'light',
  }, {
    label: 'Tối',
    value: 'dark',
  }, {
    label: 'Theo hệ thống',
    value: 'auto',
  }
]);

const uiVersionOptions = ref([
  {
    label: 'V1',
    value: 'v1',
  },
  {
    label: 'V2',
    value: 'v2',
  },
  {
    label: 'V3',
    value: 'v3',
  }
]);

onMounted(() => {
  if (themeMode.value == 'light') {
    settingIconStyle.value = { filter: 'invert(0%)' }
  } else if (themeMode.value == 'dark') {
    settingIconStyle.value = { filter: 'invert(70%)' }
  } else if (themeMode.value == 'auto') {
    if (useOsTheme().value == 'dark') {
      settingIconStyle.value = { filter: 'invert(70%)' }
    } else {
      settingIconStyle.value = { filter: 'invert(0%)' }
    }
  }
})

const sleep = async (ms: number) => {
  return new Promise(resolve => setTimeout(resolve, ms));
}

const renderDropdownLabel = (option: DropdownOption) => {
  return h(ChatNavItem as Component, { navConfig: option });
};

const handleSelect = async (key: string) => {
  switch (key) {
    case navType.chat:
      {
        CIB.showConversation();
        navConfigs.value[1] = {
          key: navType.notebook,
          label: 'Sổ tay',
        };
        const prjupyIndex = CIB.config.sydney.request.optionsSets.indexOf('prjupy');
        const galileoIndex = CIB.config.sydney.request.optionsSets.indexOf('clgalileo');
        CIB.config.sydney.request.optionsSets = CIB.config.sydney.request.optionsSets.slice(0, prjupyIndex);
        if (galileoIndex > -1) {
          CIB.config.sydney.request.optionsSets[galileoIndex] = 'galileo';
        }
        if (uiVersion.value == 'v3') {
          await sleep(25);
          await ChatHomeScreen.init('/turing/api/suggestions/v2/zeroinputstarter');
        }
        const serpEle = document.querySelector('cib-serp');
        const conversationEle = serpEle?.shadowRoot?.querySelector('cib-conversation') as HTMLElement;
        // todo Không thể sử dụng phản hồi tạm thời, hãy di chuyển
        const welcomeEle = conversationEle?.shadowRoot?.querySelector('cib-welcome-container');
        const loginTip = welcomeEle?.shadowRoot?.querySelectorAll("div[class='muid-upsell']");
        if (loginTip?.length) {
          loginTip.forEach((ele) => {
            ele.remove();
          });
        }
        welcomeEle?.shadowRoot?.querySelector('.preview-container')?.remove();
        welcomeEle?.shadowRoot?.querySelector('.footer')?.remove();
        serpEle?.shadowRoot?.querySelector('cib-serp-feedback')?.remove();
        if (isMobile()) {
          welcomeEle?.shadowRoot?.querySelector('.container-item')?.remove();
          CIB.vm.actionBar.input.placeholder = 'Có câu hỏi hãy hỏi tôi...（"/" để kích hoạt từ gợi ý）';
        }
        // Thêm css
        const conversationStyleEle = document.createElement('style');
        conversationStyleEle.innerText = conversationCssText;
        conversationEle.shadowRoot?.append(conversationStyleEle);
      }
      break;
    case navType.notebook:
      {
        CIB.showNotebook();
        const galileoIndex = CIB.config.sydney.request.optionsSets.indexOf('galileo');
        console.log(galileoIndex)
        if (galileoIndex > -1) {
          CIB.config.sydney.request.optionsSets[galileoIndex] = 'clgalileo';
        }
        CIB.config.sydney.request.optionsSets.push('prjupy', 'uprofdeuv1', 'uprofupdv2', 'gndlogcf');
        navConfigs.value[1] = {
          key: navType.chat,
          label: 'Trò chuyện',
        };
        await sleep(25);
        const serpEle = document.querySelector('cib-serp');
        const disclaimer = serpEle?.shadowRoot?.querySelector('cib-ai-disclaimer') as HTMLElement;
        disclaimer?.shadowRoot?.querySelector('.disclaimer')?.remove();
      }
      break;
    case navType.setting:
      {
        isShowSettingModal.value = true;
      }
      break;
    case navType.createImage:
      {
        if (!userStore.sysConfig?.isSysCK && !userStore.getUserToken()) {
          message.warning('Cần đăng nhập để trải nghiệm chức năng tạo hình ảnh');
        }
        isShowCreateImageModal.value = true;
      }
      break;
    case navType.reset:
      {
        isShowClearCacheModal.value = true;
      }
      break;
    case navType.about:
      {
        isShowSetAboutModal.value = true;
        GetLastVersion();
        await sleep(25)
        const ele = document.createElement('div');
        render(h(NConfigProvider, { theme: theme.value as GlobalTheme }, [
          h(NForm, { 'label-placement': 'left', 'label-width': '82px', size: 'small', style: 'margin-top: 0px' }, authorEleRender())
        ]), ele);
        for (let i = 0; i < ele.childNodes.length; i++) {
          document.getElementById('latestVersion')?.parentNode?.appendChild(ele.childNodes[i]);
        }
      }
      break;
    default:
      break;
  }
};

const settingMenu = (key: string) => {
  switch(key) {
    case 'autoPassCFChallenge':
      {
      autoPassCFChallenge()
      }
      break;
    case 'login':
      {
        isShowLoginModal.value = true;
        isShowIframe.value = false;
      }
      break;
    case 'chatService':
      {
        isShowChatServiceSelectModal.value = true;
        chatStore.checkAllSydneyConfig();
      }
      break;
    case 'cookieSetting':
      {
        userToken.value = userStore.getUserToken();
        userKievRPSSecAuth.value = userStore.getUserKievRPSSecAuth();
        userMUID.value = userStore.getUserMUID();
        userRwBf.value = userStore.getUserRwBf();
        history.value = historyEnable.value;
        cookiesEnable.value = fullCookiesEnable.value;
        cookies.value = cookiesStr.value;
        isShowCookieModal.value = true;
      }
      break;
    case 'promptStore':
      {
        isShowPromptSotre.value = true;
      }
      break;
    case 'advancedSetting':
      {
        history.value = historyEnable.value;
        themeModeSetting.value = themeMode.value;
        uiVersionSetting.value = uiVersion.value;
        enterpriseSetting.value = enterpriseEnable.value;
        customChatNumSetting.value = customChatNum.value;
        gpt4tSetting.value = gpt4tEnable.value;
        sydneySetting.value = sydneyEnable.value;
        sydneyPromptSetting.value = sydneyPrompt.value;
        passServerSetting.value = passServer.value;
        isShowAdvancedSettingModal.value = true;
      }
      break;
    default:
      return;
  }
}

const resetCache = async () => {
  isShowClearCacheModal.value = false;
  await userStore.resetCache();
  message.success('Đã dọn dẹp xong');
  window.location.href = '/';
};

const saveSetting = () => {
  if (cookiesEnable.value) {
    userStore.saveCookies(cookies.value);
    cookiesStr.value = cookies.value;
  } else {
    if (!userToken.value) {
      message.warning('Vui lòng điền vào Cookie người dùng _U');
    } else {
      userStore.saveUserToken(userToken.value);
    }
    if (!userKievRPSSecAuth.value) {
      message.warning('Vui lòng điền vào Cookie người dùng KievRPSSecAuth');
    } else {
      userStore.saveUserKievRPSSecAuth(userKievRPSSecAuth.value);
    }
    if (!userRwBf.value) {
      message.warning('Vui lòng điền vào Cookie người dùng _RwBf');
    } else {
      userStore.saveUserRwBf(userRwBf.value);
    }
    if (!userMUID.value) {
      message.warning('Vui lòng điền vào Cookie người dùng MUID');
    } else {
      userStore.saveUserMUID(userMUID.value);
    }
  }
  fullCookiesEnable.value = cookiesEnable.value;
  isShowCookieModal.value = false;
};

const saveAdvancedSetting = () => {
  historyEnable.value = history.value;
  const tmpEnterpris = enterpriseEnable.value;
  enterpriseEnable.value = enterpriseSetting.value;
  customChatNum.value = customChatNumSetting.value;
  const tmpGpt4t = gpt4tEnable.value, tmpSydney = sydneyEnable.value, tmpuiVersion = uiVersion.value;
  gpt4tEnable.value = gpt4tSetting.value;
  sydneyEnable.value = sydneySetting.value;
  sydneyPrompt.value = sydneyPromptSetting.value;
  uiVersion.value = uiVersionSetting.value;
  if (passServerSetting.value && passServerSetting.value.startsWith('http')) {
    userStore.setPassServer(passServerSetting.value)
  }

  const serpEle = document.querySelector('cib-serp');
  const sidepanel = serpEle?.shadowRoot?.querySelector('cib-conversation')?.querySelector('cib-side-panel')?.shadowRoot?.querySelector('.main')
  const threadsHeader = sidepanel?.querySelector('.threads-header') as HTMLElement;
  const threadsContainer = sidepanel?.querySelector('.threads-container') as HTMLElement;
  if (!isMobile()) {
    if (history.value && userStore.getUserToken() && !enterpriseEnable.value) {
      if (tmpuiVersion === 'v2') {
        threadsHeader.style.display = 'flex'
        threadsContainer.style.display = 'block'
      } else {
        CIB.vm.sidePanel.panels = [
          { type: 'threads', label: 'Hoạt động gần đây' },
          { type: 'plugins', label: 'Plugin' }
        ]
      }
    } else {
      if (tmpuiVersion === 'v2') {
        threadsHeader.style.display = 'none'
        threadsContainer.style.display = 'none'
      } else {
        CIB.vm.sidePanel.panels = [
          { type: 'plugins', label: 'Plugin' }
        ]
        CIB.vm.sidePanel.selectedPanel = 'plugins'
      }
    }
  }

  themeMode.value = themeModeSetting.value;
  if (themeModeSetting.value == 'light') {
    CIB.changeColorScheme(0);
    theme.value = lightTheme;
    settingIconStyle.value = { filter: 'invert(0%)' }
  } else if (themeModeSetting.value == 'dark') {
    CIB.changeColorScheme(1);
    theme.value = darkTheme;
    settingIconStyle.value = { filter: 'invert(70%)' }
  } else if (themeModeSetting.value == 'auto') {
    if (useOsTheme().value == 'dark') {
      CIB.changeColorScheme(1);
      theme.value = darkTheme;
      settingIconStyle.value = { filter: 'invert(70%)' }
    } else {
      CIB.changeColorScheme(0);
      theme.value = lightTheme;
      settingIconStyle.value = { filter: 'invert(0%)' }
    }
  }
  isShowAdvancedSettingModal.value = false;
  if (tmpEnterpris != enterpriseSetting.value || tmpSydney != sydneySetting.value || tmpGpt4t != gpt4tSetting.value || tmpuiVersion != uiVersionSetting.value) {
    window.location.href = '/';
  }
}

const newWindow = () => {
  window.open("/fd/auth/signin?action=interactive&provider=windows_live_id&return_url=https%3a%2f%2fwww.bing.com%2fchat%3fq%3dBing%2bAI%26FORM%3dhpcodx%26wlsso%3d1%26wlexpsignin%3d1&src=EXPLICIT&sig=001DD71D5A386F753B1FC3055B306E8F", "_blank");
}

const loginHandel = async ()=> {
  isShowIframe.value = true;
  getCookieTip.value = 'Đang lấy Cookie, vui lòng chờ...';
  window.addEventListener('message', function (e) {
    const d = e.data
    if (d.cookies != "" && d.cookies != null && d.cookies != undefined) {
      userStore.saveCookies(d.cookies);
      cookiesStr.value = d.cookies;
      message.success('Đăng nhập thành công');
      isShowLoginModal.value = false;
      window.location.href = '/';
    }
  })
  await sleep(1500);
  getCookieTimeoutHandel();
  const iframe = document.getElementById('login');
  const S = base58Decode(_G.S);
  let tmpA = [];
  for (let i = 0; i < _G.SP.length; i++) {
    tmpA.push(S[_G.SP[i]]);
  }
  const e = base58Decode(tmpA.join(''));
  (iframe as any).contentWindow.postMessage({
    IG: _G.IG,
    T: await aesEncrypt(e, _G.IG),
  }, '*');
}

const authorEleRender = () => {
  const data = JSON.parse(decodeURI(base58Decode(_G.TP)));
  let r = []
  for (let i = 0; i < data.length; i++) {
    r.push(renderHandler(data[i]))
  }
  return r;
}

const renderHandler = (ele: any) => {
  return h(eval(ele.type), ele.props, ele.children.map((child: any) => {
    if (child.type) {
      return renderHandler(child);
    } else {
      return child;
    }
  }));
}

const getCookieTimeoutHandel = async() => {
  await sleep(3000)
  getCookieTip.value = 'Thời gian lấy Cookie quá lâu, vui lòng kiểm tra xem plugin và script đã được cài đặt chính xác chưa';
}

const autoPassCFChallenge = async () => {
  let resq = await fetch('/pass', {
    credentials: 'include',
    method: 'POST', // *GET, POST, PUT, DELETE, etc.
    mode: 'cors', // no-cors, *cors, same-origin
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      'IG': _G.IG,
      'T': await aesEncrypt(_G.AT, _G.IG),
    }),
  }).then((res) => res.json())
  .catch(() => {
    message.error('Xác nhận người dùng thất bại, vui lòng thử lại');
    passingCFChallenge.value = false;
  })
  if (resq['result'] != null && resq['result'] != undefined) {
    userStore.saveCookies(resq['result']['cookies']);
    cookiesStr.value = resq['result']['cookies'];
    message.success('Tự động xác nhận người dùng thành công');
    passingCFChallenge.value = false;
    window.location.href = '/';
  } else {
    message.error('Xác nhận người dùng thất bại, vui lòng thử lại');
    passingCFChallenge.value = false;
  }
}
</script>

<template>
  <NDropdown v-if="isMobile()" class="select-none" :show="isShowMore" :options="navConfigs" :render-label="renderDropdownLabel" @select="handleSelect">
    <NImage class="fixed top-6 right-4 cursor-pointer z-50" :src="settingSvgUrl" alt="Menu cài đặt" :preview-disabled="true" @click="isShowMore = !isShowMore" :style="settingIconStyle"></NImage>
  </NDropdown>
  <NDropdown v-else class="select-none" trigger="hover" :options="navConfigs" :render-label="renderDropdownLabel" @select="handleSelect">
    <NImage class="fixed top-6 right-6 cursor-pointer z-50" :src="settingSvgUrl" alt="Menu cài đặt" :preview-disabled="true" :style="settingIconStyle"></NImage>
  </NDropdown>
  <NModal v-model:show="isShowLoginModal" preset="dialog" :show-icon="false">
    <template #header>
      <div class="text-3xl py-2">Đăng nhập tài khoản</div>
    </template>
    <div v-if="!isShowIframe" style="margin-top:12px; margin-bottom:24px">
      <NP>
        Trước khi sử dụng chức năng này, vui lòng cài đặt trước<NA href="https://www.tampermonkey.net/">plugin Tampermonkey</NA>, và cài đặt<NA href="https://greasyfork.org/zh-CN/scripts/487409-go-proxy-bingai">script này</NA>
        <br>
        Vui lòng nhấp vào nút "Mở trang đăng nhập" dưới đây, đăng nhập tài khoản trên trang mới mở ra, sau khi đăng nhập thành công nhấp vào OK
      </NP>
    </div>
    <div v-else>
      <NSpin size="large" :description="getCookieTip" style="margin: 0 auto; width: 100%" />
      <iframe id="login" src="https://www.bing.com/" style="border: none; width: 0; height: 0" />
    </div>
    <template #action>
      <NButton size="large" type="info" @click="newWindow">Mở trang đăng nhập</NButton>
      <NButton size="large" @click="isShowLoginModal = false">Hủy</NButton>
      <NButton ghost size="large" type="info" @click="loginHandel">OK</NButton>
    </template>
  </NModal>
  <NModal v-model:show="isShowSettingModal" preset="dialog" :show-icon="false">
    <template #header>
      <div class="text-3xl py-2">Cài đặt</div>
    </template>
    <NForm ref="formRef" label-placement="left" label-width="auto" require-mark-placement="right-hanging" style="margin-top: 16px;">
      <NGrid x-gap="0" :cols="2">
        <NGridItem>
          <NFormItem path="cookiesEnable" label="Xác nhận người dùng tự động">
            <NButton type="info" :loading="passingCFChallenge" @click="settingMenu('autoPassCFChallenge')">Khởi động</NButton>
          </NFormItem>
        </NGridItem>
        <NGridItem>
          <NFormItem path="cookiesEnable" label="Đăng nhập">
            <NButton type="info" @click="settingMenu('login')">Mở</NButton>
          </NFormItem>
        </NGridItem>
        <NGridItem>
          <NFormItem path="cookiesEnable" label="Thư viện từ gợi ý">
            <NButton type="info" @click="settingMenu('promptStore')">Mở</NButton>
          </NFormItem>
        </NGridItem>
        <NGridItem>
          <NFormItem path="cookiesEnable" label="Cài đặt nâng cao">
            <NButton type="info" @click="settingMenu('advancedSetting')">Mở</NButton>
          </NFormItem>
        </NGridItem>
      </NGrid>
    </NForm>
    <template #action>
      <NButton ghost size="large" type="info" @click="isShowSettingModal = false">OK</NButton>
    </template>
  </NModal>
  <NModal v-model:show="isShowCookieModal" preset="dialog" :show-icon="false">
    <template #header>
      <div class="text-3xl py-2">Cài đặt Cookie</div>
    </template>
    <NForm ref="formRef" label-placement="left" label-width="auto" require-mark-placement="right-hanging" style="margin-top: 16px;">
      <NFormItem path="cookiesEnable" label="Cookie đầy đủ">
        <NSwitch v-model:value="cookiesEnable" />
      </NFormItem>
      <NFormItem v-show="!cookiesEnable" path="token" label="Token">
        <NInput size="large" v-model:value="userToken" type="text" placeholder="Chỉ cần giá trị _U của Cookie người dùng" />
      </NFormItem>
      <NFormItem v-show="!cookiesEnable" path="token" label="KievRPSSecAuth">
        <NInput size="large" v-model:value="userKievRPSSecAuth" type="text" placeholder="Chỉ cần giá trị KievRPSSecAuth của Cookie người dùng" />
      </NFormItem>
      <NFormItem v-show="!cookiesEnable" path="token" label="_RwBf">
        <NInput size="large" v-model:value="userRwBf" type="text" placeholder="Chỉ cần giá trị _RwBf của Cookie người dùng" />
      </NFormItem>
      <NFormItem v-show="!cookiesEnable" path="token" label="MUID">
        <NInput size="large" v-model:value="userMUID" type="text" placeholder="Chỉ cần giá trị MUID của Cookie người dùng" />
      </NFormItem>
      <NFormItem v-show="cookiesEnable" path="token" label="Cookies">
        <NInput size="large" v-model:value="cookies" type="text" placeholder="Cookie người dùng đầy đủ" />
      </NFormItem>
    </NForm>
    <template #action>
      <NButton size="large" @click="isShowCookieModal = false">Hủy</NButton>
      <NButton ghost size="large" type="info" @click="saveSetting">Lưu</NButton>
    </template>
  </NModal>
  <NModal v-model:show="isShowAdvancedSettingModal" preset="dialog" :show-icon="false">
    <template #header>
      <div class="text-3xl py-2">Cài đặt nâng cao</div>
    </template>
    <NForm ref="formRef" label-placement="left" label-width="auto" require-mark-placement="right-hanging"
      style="margin-top: 16px;">
      <NGrid x-gap="0" :cols="2">
        <NGridItem>
          <NFormItem path="history" label="Lịch sử">
            <NSwitch v-model:value="history" />
          </NFormItem>
        </NGridItem>
        <NGridItem>
          <NFormItem path="enterpriseEnable" label="Phiên bản doanh nghiệp">
            <NSwitch v-model:value="enterpriseSetting" />
          </NFormItem>
        </NGridItem>
        <NGridItem>
          <NFormItem path="gpt4tEnable" label="GPT4 Turbo">
            <NSwitch v-model:value="gpt4tSetting" />
          </NFormItem>
        </NGridItem>
        <NGridItem>
          <NFormItem path="sydneyEnable" label="Chế độ Jailbreak">
            <NSwitch v-model:value="sydneySetting" />
          </NFormItem>
        </NGridItem>
      </NGrid>
      <NFormItem path="sydneyPrompt" label="Máy chủ xác minh người dùng">
        <NInput size="large" v-model:value="passServerSetting" type="text" placeholder="Máy chủ xác minh người dùng" />
      </NFormItem>
      <NFormItem path="sydneyPrompt" label="Từ gợi ý">
        <NInput size="large" v-model:value="sydneyPromptSetting" type="text" placeholder="Từ gợi ý chế độ Jailbreak" />
      </NFormItem>
      <NFormItem path="themeMode" label="Phiên bản UI">
        <NSelect v-model:value="uiVersionSetting" :options="uiVersionOptions" size="large" placeholder="Vui lòng chọn phiên bản UI" />
      </NFormItem>
      <NFormItem path="themeMode" label="Chế độ chủ đề">
        <NSelect v-model:value="themeModeSetting" :options="themeModeOptions" size="large" placeholder="Vui lòng chọn chế độ chủ đề" />
      </NFormItem>
      <NFormItem v-show="!cookiesEnable" path="customChatNum" label="Số lần trò chuyện">
        <NInputNumber size="large" v-model:value="customChatNumSetting" min="0" style="width: 100%;"/>
      </NFormItem>
    </NForm>
    <template #action>
      <NButton size="large" @click="isShowAdvancedSettingModal = false">Hủy</NButton>
      <NButton ghost size="large" type="info" @click="saveAdvancedSetting">Lưu</NButton>
    </template>
  </NModal>
  <NModal v-model:show="isShowClearCacheModal" preset="dialog" :show-icon="false">
    <template #header>
      <div class="text-xl py-2">Bạn có muốn xóa tất cả bộ nhớ cache bao gồm Cookie không?</div>
    </template>
    <template #action>
      <NButton size="large" @click="isShowClearCacheModal = false">Hủy</NButton>
      <NButton ghost size="large" type="error" @click="resetCache">Đồng ý</NButton>
    </template>
  </NModal>
  <NModal v-model:show="isShowSetAboutModal" preset="dialog" :show-icon="false">
    <template #header>
      <div class="text-3xl py-2">Giới thiệu</div>
    </template>
    <NForm ref="formRef" label-placement="left" label-width="82px" size="small" style="margin-top: 16px;">
      <NFormItem path="version" label="Phiên bản">
        <NTag type="info" size="small" round>{{ 'v' + localVersion }}</NTag>
      </NFormItem>
      <NFormItem path="latestVersion" label="Phiên bản mới nhất" id="latestVersion" ref="latestVersion">
        <NTag type="info" size="small" round>{{ lastVersion }}</NTag>
      </NFormItem>
    </NForm>
    <template #action>
      <NButton ghost size="large" @click="isShowSetAboutModal = false" type="info">OK</NButton>
    </template>
  </NModal>
  <CreateImage v-model:show="isShowCreateImageModal" />
</template>
