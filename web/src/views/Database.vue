<template>
  <el-table :data="tableData" stripe style="width: 100%"
            @sort-change="sortChange"
            :default-sort="{ prop: sort.name, order: sort.order }"
            table-layout="auto"
            v-loading="loading">
    <el-table-column prop="Name" label="数据库" sortable="custom" />
    <el-table-column prop="Extra.database_language" label="语言" sortable="custom"/>
    <el-table-column prop="Extra.database_linesOfCode" label="代码行数" sortable="custom"/>
    <el-table-column prop="Extra.database_cliVersion" label="CodeQL版本" sortable="custom"/>

    <el-table-column prop="Extra.database_finalised" label="构建状态" sortable="custom">
      <template #default="scope">
        <el-tooltip
            v-if="scope.row.Extra.database_finalised ==='true'"
            content="构建完成"
            placement="top"
            :hide-after="10"
        >
          <el-icon color="#7ec050" :size="20" style="margin-top: 8px"><SuccessFilled /></el-icon>
        </el-tooltip>

        <el-tooltip
            v-if="scope.row.Extra.database_finalised ==='false'"
            content="构建失败或正在构建中"
            placement="top"
            :hide-after="10"
        >
          <el-icon color="#e6c081" :size="20" style="margin-top: 8px"><QuestionFilled /></el-icon>
        </el-tooltip>
      </template>
    </el-table-column>

    <el-table-column prop="ModTime" label="修改时间" sortable="custom"
                     :formatter="(row, col, value, index)=>timeFormatter(value)"
    />

    <el-table-column fixed="right" label="">
      <template #header>
        <el-button style="float: right" :icon="Plus" @click="createData" circle/>
      </template>
      <template #default="scope">
        <el-popconfirm title="确认删除?" :hide-after="0" @confirm="deleteData(scope.row.Name)">
          <template #reference>
            <el-button :icon="Delete" circle style="float: right;margin-left: 6px"/>
          </template>
        </el-popconfirm>
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

  <el-dialog v-model="dialogFormVisible" title="上传数据库">
    <el-upload
        ref="uploader"
        class="upload-demo"
        drag
        multiple
        accept=".zip"
        :action="getBaseURL()+'/database'"
        :headers="{'X-Token':getToken()}"
        :before-upload="beforeUpload"
        :on-success="uploadSuccess"
        :on-error="uploadError"
        :on-progress="uploadProgress"
    >
      <el-icon class="el-icon--upload"><upload-filled /></el-icon>
      <div class="el-upload__text">
        Drop file here or <em>click to upload</em>
      </div>
      <template #tip>
        <div class="el-upload__tip">
          only .zip file
        </div>
      </template>
    </el-upload>
  </el-dialog>
</template>

<script setup>
import {onMounted, reactive, ref} from "vue";
import {deleteDatabases, getDatabases} from "../api/database";
import {timeFormatter} from "../utils/formatter";
import {Delete, Plus, QuestionFilled, SuccessFilled} from '@element-plus/icons-vue'
import { UploadFilled } from '@element-plus/icons-vue'
import {ElMessage} from "element-plus";
import {getBaseURL} from "../utils/request";
import {getToken} from "../utils/auth";

const emit = defineEmits(["refresh"]);


const uploader = ref(null);
const createData = () => {
  // if(uploader.value){
  //   uploader.value.clearFiles()
  // }
  dialogFormVisible.value = true
}
const deleteData = (name) => {
  console.log("delete "+name)
  deleteDatabases(name).then(response => {
    fetchData();
    ElMessage.success("删除成功")
  })
}
const dialogFormVisible = ref(false)
const beforeUpload = (file) =>{
  if(!file.name.endsWith(".zip")){
    ElMessage.error("只支持上传zip文件")
    return false
  }
  return  true
}
const uploadSuccess = (response, uploadFile, uploadFiles) => {
  if (response.err) {
    ElMessage.error(response.err)
  }else {
    ElMessage.success("数据库导入成功")
    fetchData()
  }
}
const uploadError = (error, uploadFile, uploadFiles) => {
  ElMessage.error(error.toString())
}
const uploadProgress = (evt, uploadFile, uploadFiles) => {
  if(evt.percent===100){
    ElMessage.success("上传成功，正在解压，请稍后")
  }
}


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
  