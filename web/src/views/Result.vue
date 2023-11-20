<template xmlns="http://www.w3.org/1999/html">
  <el-table :data="tableData" stripe style="width: 100%"
            @sort-change="sortChange"
            :default-sort="{ prop: sort.name, order: sort.order }"
            table-layout="auto"
            v-loading="loading"
            @filter-change="filterChange"
            :row-class-name="rowClassName_0">
    <el-table-column type="expand">
      <template #default="props">
        <div style="">
          <el-table
              :data="props.row.CodeQLSarif.Results"
              :show-header="false"
              table-layout="auto"
              :row-class-name="rowClassName">
            <el-table-column width="20px"/>
            <el-table-column type="expand">
              <template #default="scope">
                <el-table
                    :data="scope.row.CodeFlows"
                    :show-header="false"
                    table-layout="auto"
                    :row-class-name="rowClassName_S">
                  <el-table-column width="40px"/>
                  <el-table-column type="expand">
                    <template #default="scope_s">
                      <el-table
                          :data="scope_s.row.Location"
                          :show-header="false"
                          table-layout="auto"
                          :row-class-name="rowClassName_T">
                        <el-table-column width="100px"/>
                        <el-table-column>
                          <template #default="scope_t">
                            <el-tooltip
                                effect="light"
                                raw-content
                                :content="getCodeSnippet(scope_t.row,props.row.Task.ProjectLanguage)"
                                placement="right-start"
                            >
                              <el-link
                                  type="primary"
                                  :href="`https://github.com/${props.row.Task.ProjectOwner}/${props.row.Task.ProjectRepo}/blob/${props.row.Commit}/${scope_t.row.Uri}#L${scope_t.row.StartLine}-L${scope_t.row.EndLine===0?scope_t.row.StartLine:scope_t.row.EndLine}`"
                                  target="_blank"
                              >{{ scope_t.row.Message }}
                              </el-link>
                            </el-tooltip>
                          </template>
                        </el-table-column>
                        <el-table-column>
                          <template #default="scope_t">
                            <el-link
                                type="primary"
                                :href="`https://github.com/${props.row.Task.ProjectOwner}/${props.row.Task.ProjectRepo}/blob/${props.row.Commit}/${scope_t.row.Uri}#L${scope_t.row.StartLine}-L${scope_t.row.EndLine===0?scope_t.row.StartLine:scope_t.row.EndLine}`"
                                target="_blank"
                            >{{ scope_t.row.Uri }}:{{ scope_t.row.StartLine }}:{{ scope_t.row.StartColumn }}
                            </el-link>
                          </template>
                        </el-table-column>
                      </el-table>
                    </template>
                  </el-table-column>
                  <el-table-column>
                    <template #default="scope_s">
                      path
                    </template>
                  </el-table-column>
                  <el-table-column>
                    <template #default="scope_s">
                      {{ scope_s.row.Location.length }}
                    </template>
                  </el-table-column>
                </el-table>
              </template>
            </el-table-column>
            <el-table-column>
              <template #default="scope">
                <span
                    v-for="fragment in renderedMessage(scope.row.Message, scope.row.RelatedLocations, props.row.Task.ProjectLanguage)">
                  <el-tooltip
                      v-if="fragment.super"
                      effect="light"
                      raw-content
                      :content="fragment.snippet"
                      placement="top-start"
                  >
                    <el-link
                        type="primary"
                        :href="`https://github.com/${props.row.Task.ProjectOwner}/${props.row.Task.ProjectRepo}/blob/${props.row.Commit}/${fragment.uri}/#L${fragment.startLine}-L${fragment.endLine===0?fragment.startLine:fragment.endLine}`"
                        target="_blank"
                    >
                      {{ fragment.text }}
                    </el-link>
                  </el-tooltip>

                  <span v-if="!fragment.super">
                    {{ fragment.text }}
                  </span>

                </span>
              </template>
            </el-table-column>
            <el-table-column>
              <template #default="scope">
                <el-tooltip
                    effect="light"
                    raw-content
                    :content="getCodeSnippet(scope.row.Location, props.row.Task.ProjectLanguage)"
                    placement="top-start"
                >
                  <el-link
                      type="primary"
                      :href="`https://github.com/${props.row.Task.ProjectOwner}/${props.row.Task.ProjectRepo}/blob/${props.row.Commit}/${scope.row.Location.Uri}#L${scope.row.Location.StartLine}-L${scope.row.Location.EndLine===0?scope.row.Location.StartLine:scope.row.Location.EndLine}`"
                      target="_blank"
                  >{{ scope.row.Location.Uri.split('/').pop() }}:{{
                      scope.row.Location.StartLine
                    }}:{{ scope.row.Location.StartColumn }}
                  </el-link>
                </el-tooltip>
              </template>
            </el-table-column>
            <el-table-column prop="Rule"/>
          </el-table>
        </div>
      </template>
    </el-table-column>
    <el-table-column prop="FileName" label="分析结果" sortable="custom"/>
    <el-table-column prop="Task.ProjectName" label="项目"/>
    <el-table-column prop="Version" label="扫描版本" width="104px">
      <template #default="scope">
        {{
          scope.row.Task.ProjectMode === 0 ? scope.row.Version : commitVersionFormatter(scope.row.Version)
        }}
      </template>
    </el-table-column>
    <el-table-column prop="Task.ProjectSuite" label="查询套件" width="100px">
      <template #default="scope">
        <el-tag v-for="(item, index) in scope.row.Task.ProjectSuite"
                :key="index"
                type="info"
                style="margin-right: 5px">
          {{ item }}
        </el-tag>
      </template>
    </el-table-column>
<!--    <el-table-column label="漏洞数量" width="104px">-->
<!--      <template #default="scope">-->
<!--        {{ scope.row.CodeQLSarif.Results.length }}-->
<!--      </template>-->
<!--    </el-table-column>-->
    <el-table-column label="结果数量" width="104px" sortable="custom">
      <template #default="scope">
        {{ scope.row.ResultCount }}
      </template>
    </el-table-column>
    <el-table-column width="162px" prop="CreatedAt" label="创建时间" sortable="custom"
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

    <el-table-column fixed="right" label="" width="94px">
      <template #header>
        <el-tooltip
            content="全部已读"
            placement="left-start"
            :hide-after="10"
        >
          <el-button style="float: right" :icon="FolderOpened" @click="setResultIsRead(null,true)" circle/>
        </el-tooltip>
      </template>
      <template #default="scope">
        <el-tooltip
            v-if="scope.row.IsRead"
            content="标记为未读"
            placement="left-start"
            :hide-after="10"
        >
          <el-button :icon="Folder" circle @click="setResultIsRead(scope.row.ID,false)"/>
        </el-tooltip>
        <el-tooltip
            v-if="!scope.row.IsRead"
            content="标记为已读"
            placement="left-start"
            :hide-after="10"
        >
          <el-button :icon="FolderOpened" circle @click="setResultIsRead(scope.row.ID,true)"/>
        </el-tooltip>
        <el-popconfirm title="确认删除?" :hide-after="0" @confirm="deleteData(scope.row.ID)">
          <template #reference>
            <el-button :icon="Delete" circle style="margin-left: 6px"/>
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
</template>

<style>
.el-table .warning-row {
  --el-table-tr-bg-color: var(--el-color-warning-light-9);
}

.el-table .success-row {
  --el-table-tr-bg-color: var(--el-color-success-light-9);
}

.el-table .error-row {
  --el-table-tr-bg-color: var(--el-color-error-light-9);
}

.row-expand-cover .el-table__expand-column .el-icon {
  visibility: hidden;
}
</style>

<script setup>
import {onMounted, reactive, ref} from "vue";
import {deleteResult, getResults, setIsRead} from "../api/result";
import {timeFormatter} from "../utils/formatter";
import hljs from 'highlight.js'
import 'highlight.js/styles/default.min.css'
import {CircleCheck, Delete, Folder, FolderOpened, Warning} from "@element-plus/icons-vue";
import {ElMessage} from "element-plus";

const emit = defineEmits(["refresh"]);


const setResultIsRead = (id, read) => {
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

const getCodeFragment = (location) => {
  return parseLocation(location, null).fragment
}

const getCodeSnippet = (location, projectLanguage) => {
  return parseLocation(location, projectLanguage).snippet
}

const parseLocation = (location, projectLanguage) => {
  const contextStartLine = location.ContextStartLine
  const contextSnippet = location.ContextSnippet
  const startLine = location.StartLine
  const endLine = location.EndLine === 0 ? location.StartLine : location.EndLine
  const startColumn = location.StartColumn
  const endColumn = location.EndColumn

  let fragment = ''
  let snippet = ''

  const startTag = '___choccy_tag_leixiao__'
  const stopTag = '__choccy_tag_leixiao___'

  const codeLines = contextSnippet.split("\n")
  for (let i = 0; i < codeLines.length; ++i) {
    const line = contextStartLine + i
    const content = codeLines[i]

    if (startLine === endLine) {
      if (line === startLine) {
        fragment = content.slice(startColumn - 1, endColumn - 1)
        snippet += content.slice(0, startColumn - 1) + startTag
            + content.slice(startColumn - 1, endColumn - 1) + stopTag
            + content.slice(endColumn - 1) + "\n"
      } else {
        snippet += content + "\n"
      }
    } else { //startLine不等于endLine时
      if (line === startLine) {
        fragment += content.slice(startColumn - 1) + "\n"
        snippet += content.slice(0, startColumn - 1) + startTag
            + content.slice(startColumn - 1) + "\n"
      } else if (line > startLine && line < endLine) {
        fragment += content + "\n"
        snippet += content + "\n"
      } else if (line === endLine) {
        fragment += content.slice(0, endColumn - 1)
        snippet += content.slice(0, endColumn - 1) + stopTag
            + content.slice(endColumn - 1) + "\n"
      } else {
        //line小于startLine或大于endLine
        snippet += content + "\n"
      }
    }
  }
  snippet = snippet.slice(0, -1)


  if (["go", "java", "python", "html", "xml", "properties", "cpp", "swift", "yaml", "csharp", "javascript", "ruby"].indexOf(projectLanguage) > -1) {
    snippet = hljs.highlight(snippet, {language: projectLanguage}).value
  } else {
    //snippet = hljs.highlightAuto(snippet).value
    //结果数如果太长，极耗性能，前端卡死了都，还是改成固定语言
    snippet = hljs.highlight(snippet, {language: "javascript"}).value
  }

  snippet = snippet.replace(RegExp("(<[^>]+>)?" + startTag + "(<\/[^>]+>)?"), "<mark>")
  snippet = snippet.replace(RegExp("(<[^>]+>)?" + stopTag + "(<\/[^>]+>)?"), "</mark>")
  snippet = "<pre><code>" + snippet + "</pre></code>"

  return {fragment, snippet}
}

const renderedMessage = (message, locations, language) => {
  let fragments = []
  if(locations){
    let regex = "(";
    locations.forEach(function (location) {
      regex += `\\[${location.Message}\\]\\(${location.Id}\\)|`
    });
    regex = RegExp(regex.substring(0, regex.length - 1)+")","g");

    message.split(regex).forEach(function (s) {
      if(!regex.exec(s)){
        fragments.push({
          text: s
        })
      }else {
        locations.forEach(function (location) {
          if(`[${location.Message}](${location.Id})` === s){
            fragments.push({
              text: location.Message,
              super: true,
              uri: location.Uri,
              snippet: getCodeSnippet(location, language),
              startLine: location.StartLine,
              endLine: location.EndLine
            })
          }
        })
      }
    });
  }else {
    fragments.push({
      text: message
    })
  }
  return fragments
}


const rowClassName_0 = (row, index) => {
  if (row.row.CodeQLSarif.Results.length > 0) {
    return '';
  }
  return 'row-expand-cover';
}

const rowClassName = (row, index) => {
  if (row.row.CodeFlows) {
    return 'error-row';
  } else {
    return 'error-row row-expand-cover';
  }
}

const rowClassName_S = (row, index) => {
  return 'warning-row';
}

const rowClassName_T = (row, index) => {
  return 'success-row';
}


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

const sortChange = (column) => {
  sort.name = column.prop
  sort.order = column.order
  fetchData()
}


const fetchData = () => {
  loading.value = true
  getResults(paginate.currentPage, paginate.pageSize, sort.name, sort.order, JSON.stringify(filters.value)).then(response => {
    tableData.value = response["data"];
    paginate.total = response["total"];
    loading.value = false
  }).catch(err => {
    loading.value = false
  })
}

const commitVersionFormatter = (commit) => {
  return commit.substring(0, 7)
}


const filters = ref({'is_read': [false]})
const filterChange = (f) => {
  if (f.IsRead) {
    filters.value["is_read"] = f.IsRead
  }
  fetchData()
}

const deleteData = (ID) => {
  deleteResult(ID).then(response => {
    fetchData();
    ElMessage.success("删除成功")
    emit("refresh")
  })
}

onMounted(() => {
  fetchData();
})
</script>
  