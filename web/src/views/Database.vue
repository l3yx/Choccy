<template>
  <el-table :data="tableData" stripe style="width: 100%"
            @sort-change="sortChange"
            :default-sort="{ prop: sort.name, order: sort.order }"
            table-layout="auto"
            v-loading="loading">
    <el-table-column prop="Name" label="database" sortable="custom" />
    <el-table-column prop="Extra.database_language" label="language" sortable="custom"/>
    <el-table-column prop="Extra.database_linesOfCode" label="Lines of code" sortable="custom"/>
    <el-table-column prop="Extra.database_cliVersion" label="CodeQL version" sortable="custom"/>

    <el-table-column prop="Extra.database_finalised" label="Build status" sortable="custom">
      <template #default="scope">
        <el-tooltip
            v-if="scope.row.Extra.database_finalised ==='true'"
            content="Build completed"
            placement="top"
            :hide-after="10"
        >
          <el-icon color="#7ec050" :size="20" style="margin-top: 8px"><SuccessFilled /></el-icon>
        </el-tooltip>

        <el-tooltip
            v-if="scope.row.Extra.database_finalised ==='false'"
            content="Build failed or is being built"
            placement="top"
            :hide-after="10"
        >
          <el-icon color="#e6c081" :size="20" style="margin-top: 8px"><QuestionFilled /></el-icon>
        </el-tooltip>
      </template>
    </el-table-column>

    <el-table-column prop="ModTime" label="Change the time" sortable="custom"
                     :formatter="(row, col, value, index)=>timeFormatter(value)"
    />
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

<script setup>
import {onMounted, reactive, ref} from "vue";
import {getDatabases} from "../api/database";
import {timeFormatter} from "../utils/formatter";
import {QuestionFilled,SuccessFilled } from '@element-plus/icons-vue'

const emit = defineEmits(["refresh"]);

const loading = ref(true)

const tableData = ref()
const paginate = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0,
})
const sort = reactive({
  name: "ModTime",
  order: "descending"
})

const sortChange = (column)=>{
  sort.name = column.prop
  sort.order = column.order
  fetchData()
}


const fetchData = () => {
  loading.value = true
  getDatabases(paginate.currentPage, paginate.pageSize,sort.name,sort.order).then(response => {
        tableData.value = response["data"];
        paginate.total = response["total"];
    loading.value = false
  }).catch(err => {
    loading.value = false
  })
}

onMounted(() => {
  fetchData();
})
</script>
  