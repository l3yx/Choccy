<template>
  <el-table v-loading="loading"
            :data="tableData" stripe style="width: 100%"
            @sort-change="sortChange"
            :default-sort="{ prop: sort.name, order: sort.order }"
            table-layout="auto"
  >
    <el-table-column prop="Name" label="file name" sortable="custom">
      <template #default="scope">
        <el-link type="primary" @click="showQLS(scope.row)">{{ scope.row.Name }}</el-link>
      </template>
    </el-table-column>
    <el-table-column prop="Extra.suite_description" label="describe" sortable="custom"/>
    <el-table-column prop="ModTime" label="Change the time"
                     sortable="custom"
                     :formatter="(row, col, value, index)=>timeFormatter(value)"
                     />
    <el-table-column fixed="right" label="" width="106px">
      <template #header>
        <el-button style="float: right" :icon="Plus" @click="createData" circle/>
      </template>
      <template #default="scope">
        <el-popconfirm title="confirm deletion?" :hide-after="0" @confirm="deleteData(scope.row.Name)">
          <template #reference>
            <el-button :icon="Delete" circle style="float: right;margin-left: 6px"/>
          </template>
        </el-popconfirm>
        <el-button :icon="Edit" circle @click="renameData(scope.row.Name)" style="float: right"/>
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


  <el-dialog v-model="qlsFile.visible" :title="qlsFile.path" width="95%">
    <el-input
        v-model="qlsFile.content"
        :autosize="{ minRows: 3, maxRows: 10 }"
        type="textarea"
        @change="qlsFileChange"
    />

    <el-descriptions
        border
        style="margin-top: 10px"
    >
      <el-descriptions-item label="Query quantity">{{ qlsFile.queries.length }}</el-descriptions-item>
    </el-descriptions>

    <el-table v-loading="qlsFile.loading"
              :data="qlsFile.queries"
              stripe
              style="width: 100%; margin-top: 10px"
              table-layout="auto"
    >
      <el-table-column type="expand">
        <template #default="props">
          <el-descriptions
              :column="1"
              border
          >
            <el-descriptions-item label="name">{{ props.row.name }}</el-descriptions-item>
            <el-descriptions-item label="path">{{ props.row.path }}</el-descriptions-item>
            <el-descriptions-item label="description">{{ props.row.description }}</el-descriptions-item>
          </el-descriptions>
          <el-card class="box-card" shadow="never" style="margin-top: 10px">
            <pre><code v-html="hljs.highlight(props.row.content,{language:'sql'}).value"></code></pre>
          </el-card>
        </template>
      </el-table-column>
      <el-table-column prop="id" label="id" sortable/>
      <el-table-column prop="kind" label="kind" sortable/>
      <el-table-column label="tags" width="100px">
        <template #default="scope">
          <el-tag v-for="(item, index) in scope.row.tags"
                  :key="index"
                  type="info"
                  style="margin-top: 5px;margin-right: 5px">
            {{ item }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="precision" label="precision" sortable/>
      <el-table-column prop="problem.severity" label="problem.severity" sortable/>
      <el-table-column prop="security-severity" label="security-severity" sortable/>
    </el-table>
  </el-dialog>
</template>

<style>
.el-textarea__inner[readonly] {
  background: #f5f7fa;
}
</style>

<script setup>
import {onMounted, reactive, ref} from "vue";
import {getSuites, getSuiteContent, resolveSuite, saveSuiteContent,deleteSuite,createSuite,renameSuite} from "../api/suite.js";
import {ElMessage, ElMessageBox} from "element-plus";
import {timeFormatter} from "../utils/formatter";
import hljs from 'highlight.js'
import 'highlight.js/styles/default.min.css'
import {Delete, Edit, Plus} from "@element-plus/icons-vue";


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

const sortChange = (column) => {
  sort.name = column.prop
  sort.order = column.order
  fetchData()
}

const fetchData = () => {
  loading.value = true
  getSuites(paginate.currentPage, paginate.pageSize, sort.name, sort.order).then(response => {
    tableData.value = response["data"];
    paginate.total = response["total"];
    loading.value = false
  }).catch(err => {
    loading.value = false
  })
}


const qlsFile = reactive({
  visible: false,
  loading: false,
  name: "",
  path: "",
  content: "",
  queries: []
})

const qlsFileChange = () => {
  saveSuiteContent(qlsFile.name, qlsFile.content).then(response => {
    ElMessage.success("Saved successfully")
    fetchData()
    showSuiteQueries(qlsFile.path)
  });
}

const showSuiteQueries = (path) => {
  qlsFile.queries = []
  qlsFile.loading = true;
  resolveSuite(path).then(response => {
    if(response.data){
      qlsFile.queries = response.data;
    }else {
      qlsFile.queries = [];
    }
    qlsFile.loading = false;
  }).catch(err => {
    qlsFile.queries = [err];
    qlsFile.loading = false;
  });
}

const showQLS = (row) => {
  qlsFile.path = row.Path;
  qlsFile.name = row.Name;
  getSuiteContent(row.Name).then(response => {
    qlsFile.content = response.data;
    qlsFile.visible = true;
  });
  showSuiteQueries(row.Path)
}


const createData = () => {
  ElMessageBox.prompt('Please enter a file name', '', {
    confirmButtonText: 'OK',
    cancelButtonText: 'Cancel',
  }).then(({ value }) => {
    if(value && value.trim()!==""){
      createSuite(value).then(response => {
        fetchData();
        ElMessage.success("Created successfully")
      })
    }
  })
}

const deleteData = (name) => {
  deleteSuite(name).then(response => {
    fetchData();
    ElMessage.success("successfully deleted")
  })
}

const renameData = (oldName) =>{
  ElMessageBox.prompt('Please enter a file name', '', {
    confirmButtonText: 'OK',
    cancelButtonText: 'Cancel',
    inputPlaceholder: oldName
  }).then(({ value }) => {
    if(value && value.trim()!==""){
      renameSuite(oldName,value).then(response => {
        console.log(1111)
        fetchData();
        ElMessage.success("Rename successful")
      })
    }
  })

}

onMounted(() => {
  fetchData();
})
</script>
