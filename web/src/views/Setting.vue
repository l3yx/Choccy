<template>
  <el-tabs type="border-card">
    <el-tab-pane label="环境">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card class="box-card" shadow="never">
            <template #header>
              <div class="card-header">
                <el-tooltip
                    content="设置CodeQL各项配置的绝对路径或相对于Choccy二进制的相对路径"
                    placement="right"
                    :hide-after="10"
                >
                  <span>CodeQL</span>
                </el-tooltip>
              </div>
            </template>
            <el-form label-position="left" label-width="80px">
              <el-form-item label="Cli">
                <template #label>
                  <el-tooltip
                      raw-content
                      content="CodeQL引擎二进制文件路径，下载地址：https://github.com/github/codeql-cli-binaries/releases<br>设置为codeql的话将会从环境变量中寻找"
                      placement="top"
                  >
                    <span>Cli</span>
                  </el-tooltip>
                </template>
                <el-input v-model="setting.CodeQLCli" @change="(value) => settingOnchange('CodeQLCli',value)">
                  <template #append>
                    <span>{{ state.CodeQLCli_ver }}</span>
                  </template>
                </el-input>
              </el-form-item>


              <el-form-item label="Lib">
                <template #label>
                  <el-tooltip
                      raw-content
                      content="CodeQL库路径，下载地址：https://github.com/github/codeql/tags<br>设置与否codeql都会按某种规则自动寻找"
                      placement="top"
                  >
                    <span>Lib</span>
                  </el-tooltip>
                </template>
                <el-input v-model="setting.CodeQLLib" @change="saveData">
                </el-input>
              </el-form-item>

              <el-form-item label="">
                <template #label>
                  <el-tooltip
                      raw-content
                      content="自定义包路径列表，以换行分割<br>设置与否codeql都会按某种规则自动寻找"
                      placement="top"
                      :hide-after="10"
                  >
                    <span>Packs</span>
                  </el-tooltip>
                </template>
                <el-input v-model="setting.CodeQLPacks"
                          @change="saveData"
                          :autosize="{ minRows: 1, maxRows: 5 }"
                          type="textarea">
                </el-input>
              </el-form-item>

              <el-form-item label="">
                <template #label>
                  <el-tooltip
                      raw-content
                      content="查询套件路径"
                      placement="top"
                      :hide-after="10"
                  >
                    <span>Suite</span>
                  </el-tooltip>
                </template>
                <el-input v-model="setting.CodeQLSuite"
                          @change="saveData">
                </el-input>
              </el-form-item>

              <el-form-item label="">
                <template #label>
                  <el-tooltip
                      raw-content
                      content="CodeQL数据库储存路径"
                      placement="top"
                      :hide-after="10"
                  >
                    <span>Database</span>
                  </el-tooltip>
                </template>
                <el-input v-model="setting.CodeQLDatabase" @change="saveData">
                </el-input>
              </el-form-item>

              <el-form-item label="">
                <template #label>
                  <el-tooltip
                      raw-content
                      content="CodeQL分析结果储存路径"
                      placement="top"
                      :hide-after="10"
                  >
                    <span>Result</span>
                  </el-tooltip>
                </template>
                <el-input v-model="setting.CodeQLResult" @change="saveData">
                </el-input>
              </el-form-item>
            </el-form>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="box-card" shadow="never">
            <template #header>
              <div class="card-header">
                <el-tooltip
                    content="设置CodeQL等命令行工具运行时的环境变量，如Java、Go、Maven等Path变量，网络代理等<br>默认继承自系统环境变量，并可进行引用和覆盖，引用变量语法：${变量名}"
                    raw-content
                    placement="right"
                    :hide-after="10"
                >
                  <span>环境变量</span>
                </el-tooltip>
              </div>
            </template>
            <el-input
                v-model="setting.EnvStr"
                :autosize="{ minRows: 13.8, maxRows: 13.8 }"
                type="textarea"
                @change="(value) => settingOnchange('EnvStr',value)"
            />
          </el-card>
        </el-col>
      </el-row>
      <el-card class="box-card" shadow="never" style="margin-top: 20px">
        <template #header>
          <div class="card-header">
            <span>环境变量</span>
          </div>
        </template>
        <el-descriptions
            style="margin-top: 10px;"
            :column="1"
            border
            size="small"
        >
          <el-descriptions-item
              v-for="(val, key, index) in state.Env"
              :label="key.toString()"
              class-name="my-class"
              label-class-name="my-label"
          >
            {{ val }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>
    </el-tab-pane>


    <el-tab-pane label="其他">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card class="box-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>凭据</span>
              </div>
            </template>
            <el-form label-position="left" label-width="110px">
              <el-form-item label="">
                <template #label>
                  <el-tooltip
                      raw-content
                      content="用于访问本系统的口令<br>本系统有潜在的任意命令执行和文件读取功能，请务必设置强密码"
                      placement="top"
                      :hide-after="10"
                  >
                    <span>系统认证Token</span>
                  </el-tooltip>
                </template>
                <el-input v-model="setting.SystemToken" @change="saveData"></el-input>
              </el-form-item>
              <el-form-item label="">
                <template #label>
                  <el-tooltip
                      raw-content
                      content="用于版本检测，源码和数据库下载，私有仓库代码扫描等，必须设置<br>获取地址：https://github.com/settings/tokens"
                      placement="top"
                  >
                    <span>GitHub Token</span>
                  </el-tooltip>
                </template>
                <el-input v-model="setting.GithubToken" @change="saveData"/>
              </el-form-item>
            </el-form>
          </el-card>
          <el-card class="box-card" shadow="never" style="margin-top: 20px">
            <template #header>
              <div class="card-header">
                <span>扫描</span>
              </div>
            </template>
            <el-form label-position="left" label-width="170px">
              <el-form-item label="">
                <template #label>
                  <el-tooltip
                      raw-content
                      content='首次扫描Release时，需要扫描最新发布的多少个Release版本（最少1，最多10）'
                      placement="top"
                      :hide-after="10"
                  >
                    <span>首次扫描Release数量</span>
                  </el-tooltip>
                </template>
                <el-input type="number" v-model.number="setting.FirstReleaseCount" @change="saveData">
                </el-input>
              </el-form-item>
              <el-form-item label="">
                <template #label>
                  <el-tooltip
                      raw-content
                      content='30 3-6,20-23 * * * (Minutes Hours DayOfMonth Month DayOfWeek)<br>@yearly @monthly @weekly @daily @hourly<br>@every 1h30m10s<br>表达式文档：https://pkg.go.dev/github.com/robfig/cron/v3#hdr-CRON_Expression_Format'
                      placement="top"
                  >
                    <span>定时扫描Cron表达式</span>
                  </el-tooltip>
                </template>
                <el-input v-model="setting.CronTaskSpec" @change="saveData">
                  <template #append>
                    <el-tooltip
                        content='下次执行时间'
                        placement="top"
                        :hide-after="10"
                    >{{ timeFormatter(setting.CronTaskNextTime) }}
                    </el-tooltip>
                  </template>
                </el-input>
              </el-form-item>

              <el-form-item label="">
                <template #label>
                  <el-tooltip
                      raw-content
                      content='数据库分析时，需要附加的命令行选项<br>参考：https://docs.github.com/en/code-security/codeql-cli/codeql-cli-manual/database-analyze#options'
                      placement="top"
                  >
                    <span>CodeQL附加命令行选项</span>
                  </el-tooltip>
                </template>
                <el-input v-model="setting.CodeQLAnalyzeOptions" @change="saveData"/>
              </el-form-item>

              <el-form-item label="">
                <template #label>
                  <el-tooltip
                      raw-content
                      content='结束并重新运行程序时，是否恢复执行中和队列中的任务'
                      placement="top"
                      :hide-after="10"
                  >
                    <span>任务自动恢复</span>
                  </el-tooltip>
                </template>
                <el-switch v-model="setting.AutoRecoveryTask" @change="saveData"/>
              </el-form-item>
            </el-form>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="box-card" shadow="never" >
            <template #header>
              <div class="card-header">
                <span>系统</span>
              </div>
            </template>
            <el-form label-position="left" label-width="165px">
              <el-form-item label="">
                <template #label>
                  <el-tooltip
                      raw-content
                      content='访问 "GitHub项目" 页面时，获取 "最新版本" 和 "更新时间" 的最小时间间隔（不影响任务中的更新检测）<br>单位为分钟，如果设置为0，则每次刷新页面都实时获取'
                      placement="top"
                      :hide-after="10"
                  >
                    <span>项目页面更新检测间隔</span>
                  </el-tooltip>
                </template>
                <el-input type="number" v-model.number="setting.UpdateDetectionInterval" @change="saveData">
                </el-input>
              </el-form-item>
              <el-form-item label="">
                <template #label>
                  <el-tooltip
                      raw-content
                      content='发起各种HTTPS请求时，是否忽略HTTPS证书验证<br>（调试用，正常情况下不要开启该选项）'
                      placement="top"
                      :hide-after="10"
                  >
                    <span>忽略HTTPS证书验证</span>
                  </el-tooltip>
                </template>
                <el-switch v-model="setting.SkipVerifyTLS" @change="saveData"/>
              </el-form-item>
            </el-form>
          </el-card>
          <el-card class="box-card" shadow="never" style="margin-top: 20px">
            <template #header>
              <div class="card-header">
                <span>自动已读</span>
              </div>
            </template>
            <el-form label-position="left" label-width="170px">
              <el-form-item label="">
                <template #label>
                  <el-tooltip
                      raw-content
                      content='自动已读状态为完成且无扫描内容的任务'
                      placement="top"
                      :hide-after="10"
                  >
                    <span>自动已读无扫描项的任务</span>
                  </el-tooltip>
                </template>
                <el-switch v-model="setting.AutoReadEmptyTask" @change="saveData"/>
              </el-form-item>

              <el-form-item label="">
                <template #label>
                  <el-tooltip
                      raw-content
                      content='自动已读状态为完成且结果数量为0的任务'
                      placement="top"
                      :hide-after="10"
                  >
                    <span>自动已读无结果的任务</span>
                  </el-tooltip>
                </template>
                <el-switch v-model="setting.AutoReadNoResultTask" @change="saveData"/>
              </el-form-item>

              <el-form-item label="">
                <template #label>
                  <el-tooltip
                      raw-content
                      content='自动已读状态为完成的任务'
                      placement="top"
                      :hide-after="10"
                  >
                    <span>自动已读正常完成的任务</span>
                  </el-tooltip>
                </template>
                <el-switch v-model="setting.AutoReadCompletedTask" @change="saveData"/>
              </el-form-item>

              <el-form-item label="">
                <template #label>
                  <el-tooltip
                      raw-content
                      content='自动已读结果数量为0的扫描结果'
                      placement="top"
                      :hide-after="10"
                  >
                    <span>自动已读无结果的Sarif</span>
                  </el-tooltip>
                </template>
                <el-switch v-model="setting.AutoReadNoResultResult" @change="saveData"/>
              </el-form-item>

            </el-form>
          </el-card>
        </el-col>
      </el-row>
    </el-tab-pane>
  </el-tabs>
</template>

<style>
.el-card__header {
  padding-top: 10px;
  padding-bottom: 10px;
}

.my-label {
  width: 100px;
  color: #999;
  font-weight: normal;
  background: #fff;
}

.my-class {
  max-width: 295px;
  word-break: break-all;
  word-wrap: break-word;
}
</style>

<script lang="ts" setup>
import {onMounted, reactive, ref} from "vue";
import {getSetting, saveSetting, testSetting} from '../api/setting.js'
import {ElMessage} from "element-plus";
import {setToken} from "../utils/auth.js"
import {timeFormatter} from '../utils/formatter.js'


const emit = defineEmits(["refresh"]);

const state = reactive({
  CodeQLCli_ver: '',
  Env: ''
})

const setting = ref({
  CodeQLCli: '',
  CodeQLLib: '',
  CodeQLPacks: '',
  CodeQLSuite: '',
  CodeQLDatabase: '',
  CodeQLResult: '',
  EnvStr: '',

  SystemToken: '',
  GithubToken: '',


  UpdateDetectionInterval: 0,
  SkipVerifyTLS: false,

  AutoRecoveryTask: false,
  FirstReleaseCount: 0,
  CronTaskSpec: '',
  CronTaskNextTime: '',

  AutoReadEmptyTask: false,
  AutoReadNoResultTask: false,
  AutoReadCompletedTask: false,
  AutoReadNoResultResult: false,

  CodeQLAnalyzeOptions:''
})

const settingOnchange = (key, value) => {
  getSettingTest(key, value).then(_ => {
    saveData()
  })
}

const getSettingTest = (key, value) => {
  return testSetting(key, value).then(response => {
    if (key == "CodeQLCli") {
      state.CodeQLCli_ver = response.data
    } else if (key == "EnvStr") {
      state.Env = response.data
    }
  })
}

const saveData = () => {
  saveSetting(setting.value).then(response => {
    ElMessage.success("保存成功")
    setToken(setting.value.SystemToken)
    fetchData()
  })
}

const fetchData = () => {
  return getSetting().then(response => {
    setting.value = response.data;
  })
}

onMounted(() => {
  fetchData().then(_ => {
    getSettingTest("CodeQLCli", setting.value.CodeQLCli);
    getSettingTest("EnvStr", setting.value.EnvStr)
  })
})
</script>
