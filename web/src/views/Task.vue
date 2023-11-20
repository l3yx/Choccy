<template>
  <el-table :data="tableData" stripe style="width: 100%"
            @sort-change="sortChange"
            @filter-change="filterChange"
            :default-sort="{ prop: sort.name, order: sort.order }"
            table-layout="auto"
            v-loading="loading"
  >
    <el-table-column type="expand">
      <template #default="props">
        <el-input
            v-model="props.row.Logs"
            :autosize="{ minRows: 1, maxRows: 20 }"
            type="textarea"
            readonly
            style="margin-top: 5px"
        />
      </template>
    </el-table-column>
    <el-table-column prop="ProjectName" label="项目名" sortable="custom"/>
    <el-table-column prop="ProjectLanguage" label="语言" sortable="custom"/>
    <el-table-column prop="ProjectMode" :formatter="modeFormatter" label="扫描对象" sortable="custom"/>
    <el-table-column prop="Versions" label="扫描版本" width="100px">
      <template #default="scope">
        <el-tag v-for="(item, index) in scope.row.Versions"
                :key="index"
                type="info"
                style="margin-right: 5px">
          {{ scope.row.ProjectMode===0?item:commitVersionFormatter(item) }}
        </el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="ProjectSuite" label="查询套件" width="100px">
      <template #default="scope">
        <el-tag v-for="(item, index) in scope.row.ProjectSuite"
                :key="index"
                type="info"
                style="margin-right: 5px">
          {{ item }}
        </el-tag>
      </template>
    </el-table-column>

    <el-table-column prop="Stage" label="任务阶段" sortable="custom">
      <template #default="scope">
        <el-icon :size="20" style="margin-top: 8px" v-if="scope.row.Stage === 0">
          <el-tooltip content="更新检测" placement="top" :hide-after="10">
            <Loading />
          </el-tooltip>
        </el-icon>
        <el-icon  :size="20" style="margin-top: 8px" v-if="scope.row.Stage === 1">
          <el-tooltip content="资源下载" placement="top" :hide-after="10">
            <Download />
          </el-tooltip>
        </el-icon>
        <el-icon :size="20" style="margin-top: 8px" v-if="scope.row.Stage === 2">
          <el-tooltip content="数据库构建" placement="top" :hide-after="10">
            <Setting />
          </el-tooltip>
        </el-icon>
        <el-icon :size="20" style="margin-top: 8px" v-if="scope.row.Stage === 3">
          <el-tooltip content="数据库分析" placement="top" :hide-after="10">
           <Search />
          </el-tooltip>
        </el-icon>
        <span style="margin-left:10px;vertical-align:super;">{{scope.row.AnalyzedVersions.length}}/{{scope.row.Versions.length}}</span>
      </template>
    </el-table-column>

    <el-table-column prop="Status" label="任务状态" sortable="custom"
                     column-key="Status"
                     :filters="[
                      { text: '队列中', value: 0 },
                      { text: '执行中', value: 1 },
                      { text: '执行完成', value: 2 },
                      { text: '执行失败', value: -1 },
                    ]"
                     :filtered-value="filters.status"
    >
      <template #default="scope">
        <el-tooltip
            v-if="scope.row.Status ===0"
            content="队列中"
            placement="top"
            :hide-after="10"
        >
          <el-icon  color="#919398"
                    :size="20" style="margin-top: 8px">
            <RemoveFilled />
          </el-icon>
        </el-tooltip>
        <el-tooltip
            v-if="scope.row.Status ===1"
            content="执行中"
            placement="top"
            :hide-after="10"
        >
          <el-icon color="#5a9cf8"
                   :size="20" style="margin-top: 8px"><QuestionFilled /></el-icon>
        </el-tooltip>
        <el-tooltip
            v-if="scope.row.Status ===2"
            content="执行完成"
            placement="top"
            :hide-after="10"
        >
          <el-icon  color="#7ec050"
                    :size="20" style="margin-top: 8px"><SuccessFilled /></el-icon>
        </el-tooltip>
        <el-tooltip
            v-if="scope.row.Status ===-1"
            content="执行失败"
            placement="top"
            :hide-after="10"
        >
          <el-icon v-if="scope.row.Status ===-1" color="#e47470"
                   :size="20" style="margin-top: 8px"><CircleCloseFilled /></el-icon>
        </el-tooltip>
      </template>
    </el-table-column>

    <el-table-column prop="TotalResultsCount" label="结果数量" sortable="custom"/>
    <el-table-column prop="CreatedAt" label="创建时间" sortable="custom"
                     :formatter="(row, col, value, index)=>timeFormatter(value)"/>

    <el-table-column
        width="66px"
        label="查阅"
        prop="IsRead"
        column-key="IsRead"
        :filters="[
          { text: '已读', value: true },
          { text: '未读', value: false },
        ]"
        :filtered-value="filters.is_read">
      <template #default="scope">
        <el-icon v-if="scope.row.IsRead" :size="20" color="#a3d280" style="margin-top: 8px">
          <CircleCheck/>
        </el-icon>
        <el-icon v-if="!scope.row.IsRead" :size="20" style="margin-top: 8px">
          <Warning/>
        </el-icon>
      </template>
    </el-table-column>
    <el-table-column fixed="right" label="" width="50px">
      <template #header>
        <el-tooltip
            content="全部已读"
            placement="left-start"
            :hide-after="10"
        >
          <el-button style="float: right" :icon="FolderOpened" @click="setTaskIsRead(null,true)" circle/>
        </el-tooltip>
      </template>
      <template #default="scope">
        <el-tooltip
            v-if="scope.row.IsRead"
            content="标记为未读"
            placement="left-start"
            :hide-after="10"
        >
          <el-button :icon="Folder" circle @click="setTaskIsRead(scope.row.ID,false)"/>
        </el-tooltip>
        <el-tooltip
            v-if="!scope.row.IsRead"
            content="标记为已读"
            placement="left-start"
            :hide-after="10"
        >
          <el-button :icon="FolderOpened" circle @click="setTaskIsRead(scope.row.ID,true)"/>
        </el-tooltip>
      </template>
    </el-table-column>
  </el-table>

  <el-pagination
      style="margin-top: 20px"
      v-model:current-page="paginate.currentPage"
      v-model:page-size="paginate.pageSize"
      :page-sizes="[1, 5, 10, 15, 20, 50]"
      layout="total, sizes, prev, pager, next"
      v-model:total="paginate.total"
      @size-change="fetchData"
      @current-change="fetchData"
  />
</template>

<style>
.el-textarea__inner[readonly] {
  background: #f5f7fa;
}
.el-step__icon{
  width: 20px;
  height: 20px;
}
.el-step__title{
  font-size: 14px;
}
</style>

<script setup>
import {onMounted, reactive, ref} from "vue";
import {getTasks,setIsRead} from "../api/task.js";
import {timeFormatter} from "../utils/formatter";
import {
  RemoveFilled,
  QuestionFilled,
  SuccessFilled,
  CircleCloseFilled,
  Download,
  Setting,
  Search,Loading,
  FolderOpened, Folder, Warning, CircleCheck
} from '@element-plus/icons-vue'

const emit = defineEmits(["refresh"]);
const loading = ref(true)

const tableData = ref()
const paginate = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0,
})
const sort = reactive({
  name: "CreatedAt",
  order: "descending"
})
const filters = ref({
      'status':[0,1,2,-1],
      'is_read': [false]
})

const sortChange = (column) => {
  sort.name = column.prop
  sort.order = column.order
  fetchData()
}

const filterChange = (f) => {
  if (f.Status) {
    filters.value["status"] = f.Status
  }
  if (f.IsRead) {
    filters.value["is_read"] = f.IsRead
  }
  fetchData()
}

const fetchData = () => {
  loading.value = true
  getTasks(paginate.currentPage, paginate.pageSize, sort.name, sort.order, JSON.stringify(filters.value)).then(response => {
    tableData.value = response["data"];
    paginate.total = response["total"];
    loading.value = false
  }).catch(err => {
    loading.value = false
  })
}

const modeFormatter = (row, col, value, index) => {
  if (value === 0) {
    return "Release";
  } else if (value === 1) {
    return "原有数据库";
  }
  return value;
}

const commitVersionFormatter = (commit) =>{
  return commit.substring(0, 7)
}

const setTaskIsRead = (id, read) => {
  let idList = []
  if (id !== null) {
    idList.push(id)
  } else {
    tableData.value.forEach(function (item) {
      idList.push(item.ID)
    });
  }
  setIsRead(idList, read).then(response => {
    fetchData()
    emit("refresh")
  })
}

onMounted(() => {
  fetchData();
})
</script>
