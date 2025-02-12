export default {
  route: {
    dashboard: '首页',
    pipeline: '流水线',
    documentation: '文档',
    guide: '引导页',
    permission: '权限测试页',
    rolePermission: '角色权限',
    pagePermission: '页面权限',
    directivePermission: '指令权限',
    icons: '图标',
    components: '组件',
    tinymce: '富文本编辑器',
    markdown: 'Markdown',
    jsonEditor: 'JSON 编辑器',
    dndList: '列表拖拽',
    splitPane: 'Splitpane',
    avatarUpload: '头像上传',
    dropzone: 'Dropzone',
    sticky: 'Sticky',
    countTo: 'Count To',
    componentMixin: '小组件',
    backToTop: '返回顶部',
    dragDialog: '拖拽 Dialog',
    dragSelect: '拖拽 Select',
    dragKanban: '可拖拽看板',
    charts: '图表',
    keyboardChart: '键盘图表',
    lineChart: '折线图',
    mixChart: '混合图表',
    example: '综合实例',
    nested: '路由嵌套',
    menu1: '菜单1',
    'menu1-1': '菜单 1-1',
    'menu1-2': '菜单 1-2',
    'menu1-2-1': '菜单 1-2-1',
    'menu1-2-2': '菜单 1-2-2',
    'menu1-3': '菜单 1-3',
    menu2: '菜单 2',
    Table: 'Table',
    dynamicTable: '动态 Table',
    dragTable: '拖拽 Table',
    inlineEditTable: 'Table 内编辑',
    complexTable: '综合 Table',
    tab: 'Tab',
    form: '表单',
    createArticle: '创建文章',
    editArticle: '编辑文章',
    articleList: '文章列表',
    errorPages: '错误页面',
    page401: '401',
    page404: '404',
    errorLog: '错误日志',
    excel: 'Excel',
    exportExcel: '导出 Excel',
    selectExcel: '导出 已选择项',
    mergeHeader: '导出 多级表头',
    uploadExcel: '上传 Excel',
    zip: 'Zip',
    pdf: 'PDF',
    exportZip: 'Export Zip',
    theme: '换肤',
    clipboardDemo: 'Clipboard',
    i18n: '国际化',
    externalLink: '外链',
    profile: '个人中心',
    pipelineTable: '流水线列表',
    pipelineDetail: '流水线详情',
    pipelineMonitor: '流水线监控',
    nodeTable: '节点列表',
    cluster: '集群'
  },
  navbar: {
    dashboard: '首页',
    github: '项目地址',
    logOut: '退出登录',
    profile: '个人中心',
    theme: '换肤',
    size: '布局大小'
  },
  login: {
    title: '系统登录',
    logIn: '登录',
    authType: '认证方式',
    username: '账号',
    password: '密码',
    any: '随便填',
    thirdparty: '第三方登录',
    thirdpartyTips: '本地不能模拟，请结合自己业务进行模拟！！！'
  },
  documentation: {
    documentation: '文档',
    github: 'Github 地址'
  },
  permission: {
    addRole: '新增角色',
    editPermission: '编辑权限',
    roles: '你的权限',
    switchRoles: '切换权限',
    tips: '在某些情况下，不适合使用 v-permission。例如：Element-UI 的 el-tab 或 el-table-column 以及其它动态渲染 dom 的场景。你只能通过手动设置 v-if 来实现。',
    delete: '删除',
    confirm: '确定',
    cancel: '取消'
  },
  guide: {
    description: '引导页对于一些第一次进入项目的人很有用，你可以简单介绍下项目的功能。本 Demo 是基于',
    button: '打开引导'
  },
  components: {
    documentation: '文档',
    tinymceTips: '富文本是管理后台一个核心的功能，但同时又是一个有很多坑的地方。在选择富文本的过程中我也走了不少的弯路，市面上常见的富文本都基本用过了，最终权衡了一下选择了Tinymce。更详细的富文本比较和介绍见',
    dropzoneTips: '由于我司业务有特殊需求，而且要传七牛 所以没用第三方，选择了自己封装。代码非常的简单，具体代码你可以在这里看到 @/components/Dropzone',
    stickyTips: '当页面滚动到预设的位置会吸附在顶部',
    backToTopTips1: '页面滚动到指定位置会在右下角出现返回顶部按钮',
    backToTopTips2: '可自定义按钮的样式、show/hide、出现的高度、返回的位置 如需文字提示，可在外部使用Element的el-tooltip元素',
    imageUploadTips: '由于我在使用时它只有vue@1版本，而且和mockjs不兼容，所以自己改造了一下，如果大家要使用的话，优先还是使用官方版本。'
  },
  table: {
    dynamicTips1: '固定表头, 按照表头顺序排序',
    dynamicTips2: '不固定表头, 按照点击顺序排序',
    dragTips1: '默认顺序',
    dragTips2: '拖拽后顺序',
    title: '标题',
    importance: '重要性',
    type: '类型',
    remark: '点评',
    search: '搜索',
    add: '添加',
    export: '导出',
    reviewer: '审核人',
    id: '序号',
    date: '时间',
    author: '作者',
    readings: '阅读数',
    status: '状态',
    actions: '操作',
    edit: '编辑',
    publish: '发布',
    draft: '草稿',
    delete: '删除',
    cancel: '取 消',
    confirm: '确 定'
  },
  example: {
    warning: '创建和编辑页面是不能被 keep-alive 缓存的，因为keep-alive 的 include 目前不支持根据路由来缓存，所以目前都是基于 component name 来进行缓存的。如果你想类似的实现缓存效果，可以使用 localStorage 等浏览器缓存方案。或者不要使用 keep-alive 的 include，直接缓存所有页面。详情见'
  },
  errorLog: {
    tips: '请点击右上角bug小图标',
    description: '现在的管理后台基本都是spa的形式了，它增强了用户体验，但同时也会增加页面出问题的可能性，可能一个小小的疏忽就导致整个页面的死锁。好在 Vue 官网提供了一个方法来捕获处理异常，你可以在其中进行错误处理或者异常上报。',
    documentation: '文档介绍'
  },
  excel: {
    export: '导出',
    selectedExport: '导出已选择项',
    placeholder: '请输入文件名(默认excel-list)'
  },
  zip: {
    export: '导出',
    placeholder: '请输入文件名(默认file)'
  },
  pdf: {
    tips: '这里使用   window.print() 来实现下载pdf的功能'
  },
  theme: {
    change: '换肤',
    documentation: '换肤文档',
    tips: 'Tips: 它区别于 navbar 上的 theme-pick, 是两种不同的换肤方法，各自有不同的应用场景，具体请参考文档。'
  },
  tagsView: {
    refresh: '刷新',
    close: '关闭',
    closeOthers: '关闭其它',
    closeAll: '关闭所有'
  },
  settings: {
    title: '系统布局配置',
    theme: '主题色',
    tagsView: '开启 Tags-View',
    fixedHeader: '固定 Header',
    sidebarLogo: '侧边栏 Logo'
  },
  pipeline: {
    name: '流水线名称',
    aliasName: '流水线别名',
    status: '状态',
    statusValues: {
      started: '已启动',
      starting: '启动中',
      stopped: '已停止',
      stopping: '停止中'
    },
    statusOptions: [
      {
        key: 'run',
        value: '运行'
      },
      {
        key: 'stop',
        value: '停止'
      }
    ],
    infoStatusMap: {
      run: {
        key: 'run',
        value: '运行',
        type: 'success'
      },
      stop: {
        key: 'stop',
        value: '停止',
        type: 'info'
      },
      stopping: {
        key: 'stopping',
        value: '停止中',
        type: 'warning'
      },
      scheduling: {
        key: 'scheduling',
        value: '调度中',
        type: 'warning'
      },
      scheduled: {
        key: 'scheduled',
        value: '调度完成',
        type: ''
      },
      deleting: {
        key: 'deleting',
        value: '删除中',
        type: 'danger'
      }
    },
    createTime: '创建时间',
    mysqlAddress: 'Mysql地址',
    mysqlPort: 'Mysql端口',
    mysqlUser: 'Mysql用户名',
    mysqlPassword: 'Mysql密码',
    mysqlServerId: 'Mysql从服务器ID',
    mysqlFlavor: 'Mysql类型',
    mysqlFlavorOptions: [
      {
        key: 'MySQL',
        value: 'MySQL'
      },
      {
        key: 'MariaDB',
        value: 'MariaDB'
      }
    ],
    mysqlMode: 'Mysql binlog模式',
    mysqlModeOptions: [
      {
        key: 'position',
        value: '普通'
      },
      {
        key: 'gtid',
        value: 'GTID'
      }
    ],
    output: {
      sender: {
        name: '输出名称',
        type: '输出类型',
        typeOptions: [
          {
            key: 'stdout',
            value: '标准输出'
          },
          {
            key: 'kafka',
            value: 'Kafka'
          },
          {
            key: 'http',
            value: 'HTTP'
          },
          {
            key: 'rabbitMQ',
            value: 'RabbitMQ'
          },
          {
            key: 'redis',
            value: 'Redis'
          },
          {
            key: 'rocketMQ',
            value: 'RocketMQ'
          }
          // {
          //   key: 'elastic',
          //   value: 'Elasticsearch'
          // }
        ]
      }
    },
    bindNode: '绑定节点',
    remark: '备注',
    runNode: '运行节点',
    filter: {
      type: '类型',
      rule: '规则'
    },
    filterValues: {
      white: '白名单',
      black: '黑名单'
    }
  },
  node: {
    createTime: '首次加入时间',
    name: '节点名称',
    leader: '主节点',
    ip: '节点当前IP',
    status: '节点状态',
    version: '节点版本',
    roles: '角色',
    roleMap: {
      leader: '主节点',
      follower: '从节点',
      master: '候选人',
      admin: '管理后台',
      worker: '工作节点'
    }
  },
  capacity: {
    cpuCores: 'CPU核心数',
    memory: '内存容量',
    cpuUsage: 'CPU使用率',
    memoryUsage: '内存使用率'
  },
  pipeline_table: {
    select: {
      pleaseSelect: '请选择'
    },
    date: {
      pickDate: '请选择日期'
    },
    input: {
      pleaseInput: '请输入'
    },
    filter: {
      name: '过滤器',
      delete: '删除',
      addFilter: '新增过滤器',
      whiteList: '白名单',
      blackList: '黑名单',
      place: '输入database.table或者database'
    },
    kafka: {
      brokers: '多个broker使用逗号分隔',
      acksOptions: [
        {
          key: 0,
          label: 'NoResponse'
        }
      ]
    },
    pipeline: {
      name: '流水线名称，全局唯一',
      aliasName: '流水线别名，随意取，便于理解',
      run: '运行'
    },
    http: {
      retries: '重试次数'
    },
    rabbit: {
      tips: '配置rabbitMQ的exchange，binlogo会向该exchange发送message，routing key是数据库加数据表，形如database.table, 用户使用时需要创建相应的rabbitMQ队列'
    },
    redis: {
      tips: 'Binlogo会向redis的list发送数据，使用的命令是RPUSH,从list右方插入数据'
    },
    rocket: {
      tips: '因为一般使用RocketMQ的都是直接使用阿里云的云服务, 所以这里直接提供了阿里云的配置,参考地址'
    },
    tips: {
      fix_pos_newest: '因为mysql删除了binlog文件导致无法同步时，是否使用最新的binlog位置信息?'
    },
    actions: '操作',
    status: '状态',
    edit: '编辑',
    delete: '删除',
    run: '运行',
    stop: '停止',
    name: '流水线名称',
    server_id: '如果为空后台会自动生成',
    // sortOptions: [{ label: 'Time Ascending', key: '+id' }, { label: 'Time Descending', key: '-id' }]
    sortOptions: [{ label: '升序', key: '+id' }, { label: '降序', key: '-id' }],
    statusOptions: [{ label: '运行', key: 'run' }, { label: '停止', key: 'stop' }]
  },
  pipe_detail: {
    pipeline: '流水线',
    remark: '备注',
    detail: '详情',
    runStatus: '运行状态',
    instance: '实例',
    govern: '治理',
    filter: '过滤器',
    event: '事件',
    filterTip: '输入表名，验证当前过滤规则下该数据库或数据表能否通过，表名格式为数据库名称加表名，例如 test_db.test_tbl',
    validRule: '待验证的数据库名或数据表名',
    button_valid: '验证',
    validTipTrue: ' 将会被过滤',
    validTipFalse: ' 不会被过滤',
    positionTip: '手动选择一个binlog同步的位置。用于故障恢复或特殊的业务需求。',
    eventTip: '显示最近的20条事件记录',
    event_table: {
      first_time: '首次记录时间',
      last_time: '最后记录时间',
      type: '类型',
      count: '聚合次数',
      node_name: '上报节点',
      node_ip: '节点IP',
      message: '事件内容'
    }
  },
  node_table: {
    statusMap: {
      ready: {
        yes: '就绪',
        no: '未就绪'
      },
      network_unavailable: {
        yes: '网络不可达'
      }
    },
    name: '节点名称',
    ready: '状态',
    sortOptions: [{ label: '升序', key: '+id' }, { label: '降序', key: '-id' }],
    readyOptions: [{ label: '就绪', key: 'yes' }, { label: '未就绪', key: 'no' }]
  },
  instance: {
    create_time: '启动时间',
    pipeline_name: '流水线',
    node_name: '节点名称'
  },
  register: {
    create_time: '注册时间',
    name: '节点名称',
    ip: '节点IP',
    version: '版本'
  },
  cluster: {
    title: '集群信息',
    name: '集群名称',
    tabMap: {
      election: {
        key: 'election',
        label: '选举节点'
      },
      register: {
        key: 'register',
        label: '注册节点'
      },
      instance: {
        key: 'instance',
        label: '实例'
      }
    },
    instanceTip: '显示当前正在运行的流水线实例.',
    registerTip: '显示当前所有的已经注册（在线）的节点.',
    electionTip: '显示主节点的竞选情况.'
  },
  global: {
    submit: '提交'
  }
}
