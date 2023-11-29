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
    <el-table-column prop="Project Name" label="Item name" sortable="custom"/>
    <el-table-column prop="Project Language" label="language" sortable="custom"/>
    <el-table-column prop="Project Mode" :formatter="modeFormatter" label="Scan object" sortable="custom"/>
    <el-table-column prop="Versions" label="Scanned version" width="100px">
      <template #default="scope">
        <el-tag v-for="(item, index) in scope.row.Versions"
                :key="index"
                type="info"
                style="margin-right: 5px">
          {{ scope.row.ProjectMode===0?item:commitVersionFormatter(item) }}
        </el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="Project Suite" label="query suite" width="100px">
      <template #default="scope">
        <el-tag v-for="(item, index) in scope.row.ProjectSuite"
                :key="index"
                type="info"
                style="margin-right: 5px">
          {{ item }}
        </el-tag>
      </template>
    </el-table-column>

    <el-table-column prop="Stage" label="mission phase" sortable="custom">
      <template #default="scope">
        <el-icon :size="20" style="margin-top: 8px" v-if="scope.row.Stage === 0">
          <el-tooltip content="Update detection" placement="top" :hide-after="10">
            <Loading />
          </el-tooltip>
        </el-icon>
        <el-icon  :size="20" style="margin-top: 8px" v-if="scope.row.Stage === 1">
          <el-tooltip content="Download" placement="top" :hide-after="10">
            <Download />
          </el-tooltip>
        </el-icon>
        <el-icon :size="20" style="margin-top: 8px" v-if="scope.row.Stage === 2">
          <el-tooltip content="Database construction" placement="top" :hide-after="10">
            <Setting />
          </el-tooltip>
        </el-icon>
        <el-icon :size="20" style="margin-top: 8px" v-if="scope.row.Stage === 3">
          <el-tooltip content="Database analysis" placement="top" :hide-after="10">
           <Search />
          </el-tooltip>
        </el-icon>
        <span style="margin-left:10px;vertical-align:super;">{{scope.row.AnalyzedVersions.length}}/{{scope.row.Versions.length}}</span>
      </template>
    </el-table-column>

    <el-table-column prop="Status" label="Task status" sortable="custom"
                     column-key="Status"
                     :filters="[
                      { text: 'in queue', value: 0 },
                      { text: 'Executing', value: 1 },
                      { text: 'Execution completed', value: 2 },
                      { text: 'Execution failed', value: -1 },
                    ]"
                     :filtered-value="filters.status"
    >
      <template #default="scope">
        <el-tooltip
            v-if="scope.row.Status ===0"
            content="in queue"
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
            content="Executing"
            placement="top"
            :hide-after="10"
        >
          <el-icon color="#5a9cf8"
                   :size="20" style="margin-top: 8px"><QuestionFilled /></el-icon>
        </el-tooltip>
        <el-tooltip
            v-if="scope.row.Status ===2"
            content="Execution completed"
            placement="top"
            :hide-after="10"
        >
          <el-icon  color="#7ec050"
                    :size="20" style="margin-top: 8px"><SuccessFilled /></el-icon>
        </el-tooltip>
        <el-tooltip
            v-if="scope.row.Status ===-1"
            content="Execution failed"
            placement="top"
            :hide-after="10"
        >
          <el-icon v-if="scope.row.Status ===-1" color="#e47470"
                   :size="20" style="margin-top: 8px"><CircleCloseFilled /></el-icon>
        </el-tooltip>
      </template>
    </el-table-column>

    <el-table-column prop="TotalResultsCount" label="number of results" sortable="custom"/>
    <el-table-column prop="CreatedAt" label="creation time" sortable="custom"
                     :formatter="(row, col, value, index)=>timeFormatter(value)"/>

    <el-table-column
        width="66px"
        label="Check"
        prop="IsRead"
        column-key="IsRead"
        :filters="[
          { text: 'Have read', value: true },
          { text: 'unread', value: false },
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
            content="All read"
            placement="left-start"
            :hide-after="10"
        >
          <el-button style="float: right" :icon="FolderOpened" @click="setTaskIsRead(null,true)" circle/>
        </el-tooltip>
      </template>
      <template #default="scope">
        <el-tooltip
            v-if="scope.row.IsRead"
            content="Mark as unread"
            placement="left-start"
            :hide-after="10"
        >
          <el-button :icon="Folder" circle @click="setTaskIsRead(scope.row.ID,false)"/>
        </el-tooltip>
        <el-tooltip
            v-if="!scope.row.IsRead"
            content="Mark as read"
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
    return "Original database";
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
