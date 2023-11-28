<template>
  <div>
    <el-container>
      <el-header>
        <el-menu :default-active="router.currentRoute.value.path" class="el-menu-demo" mode="horizontal"
                 :ellipsis="false"
                 @select="handleSelect">
          <el-menu-item index="0">
            <el-icon :size="56">
              <ChoccyIcon/>
            </el-icon>
          </el-menu-item>
          <div class="flex-grow"/>
          <el-menu-item index="/project">GitHub project</el-menu-item>
          <el-menu-item index="/pack">Query package</el-menu-item>
          <el-menu-item index="/suite">query suite</el-menu-item>
          <el-menu-item index="/database">database</el-menu-item>
          <el-menu-item index="/task">
            <el-badge :value="unread.task" :hidden="unread.task===0" class="item">Task</el-badge>
          </el-menu-item>
          <el-menu-item index="/result">
            <el-badge :value="unread.result" :hidden="unread.result===0" class="item">Analyze results</el-badge>
          </el-menu-item>
          <el-menu-item index="/setting">set up</el-menu-item>
        </el-menu>
      </el-header>
      <el-main style="">
        <RouterView @refresh="refresh"/>
      </el-main>
    </el-container>
  </div>
</template>

<script lang="ts" setup>
import ChoccyIcon from './components/IconChoccy.vue'
import router from "./router/index.js";
import {onMounted, onUpdated, reactive} from "vue";
import {getNotifications} from './api/notification.js'
import {setToken} from './utils/auth.js'
import {ElMessageBox, ElNotification} from "element-plus";
import {getResultUnread} from './api/result.js'
import {getTaskUnread} from './api/task.js'

const fetchNotifications = () => {
  getNotifications().then(response => {
    if (response.data.notifications.length>0) {
      ElNotification({
        title: 'hint',
        position: 'top-left',
        showClose: true,
        dangerouslyUseHTMLString: true,
        message: parseNotifications(response.data.notifications),
        duration: 0,
      })
    }
  }).catch(err => {
    if (err === "Unauthorized") {
      openAuth()
    }
  })
}

const parseNotifications = (notifications) => {
  let str = "<ul style='padding-left: 20px'>"
  notifications.forEach(function (notification) {
    str += `<li><span style='color: teal'>${notification}</span></li>`
  });
  str += "</ul>"
  return str
}

const openAuth = () => {
  ElMessageBox.prompt('Please enter system Token', '', {
    confirmButtonText: 'OK',
    showClose: false,
    showCancelButton: false,
    closeOnClickModal: false,
    closeOnPressEscape: false,
    closeOnHashChange: false,
  }).then(({value}) => {
    setToken(value)
    location.reload()
  })
}

const unread = reactive({
  task:0,
  result:0
})
const fetchUnread = () => {
  getResultUnread().then(response=>{
    unread.result = response["count"]
  })
  getTaskUnread().then(response=>{
    unread.task = response["count"]
  })
}

const refresh = () => {
  fetchNotifications()
  fetchUnread()
}

function handleSelect(key: string, keyPath: string[]) {
  if (key === "0") {
    window.open("https://github.com/l3yx/choccy", '_blank');
  } else {
    router.push(key)
  }
  refresh()
}

onMounted(() => {
  refresh()
})

onUpdated(()=>{
  console.log(11)
})
</script>


<style scoped>
.flex-grow {
  flex-grow: 1;
}

body {
  margin: 0px;
}

pre {
  margin: 0px;
}

.el-menu {
  --el-menu-item-height: 30px;
  height: 45px;
}
.el-sub-menu .el-menu-item{
  height: 45px;
}
</style>
