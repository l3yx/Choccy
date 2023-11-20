<template>
  <el-tree v-loading="loading" :props="props" :load="loadNode" lazy node-key="id"
           :default-expanded-keys="expandedKeys"
           @node-click="nodeClick">
    <template #default="{ node, data }">
        <span class="custom-tree-node">
          <el-tooltip
              raw-content
              :disabled="!(node.level === 3)"
              :content="node.data.full"
              placement="right"
              :hide-after="10"
              :offset="12"
          >
            <el-text>
              <el-icon>
                <Document v-if="node.level===4"/>
                <FolderOpened v-if="node.level===3"/>
                <Files v-if="node.level===2"/>
                <OfficeBuilding v-if="node.level===1"/>
              </el-icon>

              <span v-if="node.level===4">{{node.data.ql.replace(node.data.path+"/","")}}</span>
              <span v-if="node.level!=4">{{ node.label }}</span>
            </el-text>
          </el-tooltip
          >
        </span>
    </template>
  </el-tree>

  <el-dialog v-model="qlFile.visible" :title="qlFile.title" width="80%">
    <el-card class="box-card" shadow="never" style="margin-top: 10px">
      <pre><code v-html="hljs.highlight(qlFile.content,{language:'sql'}).value"></code></pre>
    </el-card>
  </el-dialog>
</template>

<style>
.el-textarea__inner[readonly] {
  background: #f5f7fa;
}
</style>

<script lang="ts" setup>
import {getQueries, getQueryContent} from '../api/query.js'
import {onMounted, reactive, ref} from "vue";
import {Document, Files, FolderOpened,OfficeBuilding} from '@element-plus/icons-vue'
import type Node from 'element-plus/es/components/tree/src/model/node'
import hljs from 'highlight.js'
import 'highlight.js/styles/default.min.css'

interface Tree {
  name: string
  leaf?: boolean
}

const props = {
  label: 'name',
  isLeaf: 'leaf',
}

const emit = defineEmits(["refresh"]);

const loading = ref(true)

const qlFile = reactive({
  visible: false,
  title: "",
  content: ""
})

const expandedKeys = ref([])


const loadNode = (node: Node, resolve: (data: Tree[]) => void) => {
  if (node.level === 0) {
    loading.value = true
    fetchData().then(_ => {
      const scopes = Object.keys(data.value)
      const items = []
      scopes.forEach(element => {
        if(element !="/" && element!="codeql"){
          expandedKeys.value.push(element)
        }
        items.push({
          id: element,
          scope: element,
          name: element,
          leaf: false
        })
      });
      loading.value = false
      return resolve(items)
    }).catch(err => {
      loading.value = false
    })
  } else if (node.level === 1) {
    const packs = Object.keys(data.value[node.data.scope])
    const items = []
    packs.forEach(element => {
      items.push({
        scope: node.data.scope,
        pack: element,
        name: element,
        leaf: false
      })
    });
    return resolve(items)
  } else if (node.level === 2) {
    const paths = data.value[node.data.scope][node.data.pack]
    const items = []
    paths.forEach(element => {
      items.push({
        scope: node.data.scope,
        pack: node.data.pack,
        path: element,
        full: element,
        name: element.split('/').pop(),
        leaf: false
      })
    });
    return resolve(items)
  } else if (node.level === 3) {
    getQueries(node.data.path).then(response => {
      const items = []
      response.data.forEach(element => {
        items.push({
          scope: node.data.scope,
          pack: node.data.pack,
          path: node.data.path,
          ql: element,
          full: element,
          name: element.split('/').pop(),
          leaf: true
        })
      });
      return resolve(items)
    }).catch(err => {
      return resolve([{"name": err}])
    })
  }
}

const nodeClick = (o, node, tree, event) => {
  if (node.level === 4) {
    getQueryContent(node.data.ql).then(response => {
      qlFile.content = response.data
      qlFile.title = node.data.ql
      qlFile.visible = true
    })
  }
}

const data = ref({})

const fetchData = () => {
  return getQueries("").then(response => {
    data.value = response.data
  })
}

onMounted(() => {
  //fetchData();
})
</script>
