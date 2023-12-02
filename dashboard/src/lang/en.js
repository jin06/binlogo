export default {
  route: {
    dashboard: 'Dashboard',
    pipeline: 'Pipeline',
    documentation: 'Documentation',
    guide: 'Guide',
    permission: 'Permission',
    pagePermission: 'Page Permission',
    rolePermission: 'Role Permission',
    directivePermission: 'Directive Permission',
    icons: 'Icons',
    components: 'Components',
    tinymce: 'Tinymce',
    markdown: 'Markdown',
    jsonEditor: 'JSON Editor',
    dndList: 'Dnd List',
    splitPane: 'SplitPane',
    avatarUpload: 'Avatar Upload',
    dropzone: 'Dropzone',
    sticky: 'Sticky',
    countTo: 'Count To',
    componentMixin: 'Mixin',
    backToTop: 'Back To Top',
    dragDialog: 'Drag Dialog',
    dragSelect: 'Drag Select',
    dragKanban: 'Drag Kanban',
    charts: 'Charts',
    keyboardChart: 'Keyboard Chart',
    lineChart: 'Line Chart',
    mixChart: 'Mix Chart',
    example: 'Example',
    nested: 'Nested Routes',
    menu1: 'Menu 1',
    'menu1-1': 'Menu 1-1',
    'menu1-2': 'Menu 1-2',
    'menu1-2-1': 'Menu 1-2-1',
    'menu1-2-2': 'Menu 1-2-2',
    'menu1-3': 'Menu 1-3',
    menu2: 'Menu 2',
    Table: 'Table',
    dynamicTable: 'Dynamic Table',
    dragTable: 'Drag Table',
    inlineEditTable: 'Inline Edit',
    complexTable: 'Complex Table',
    tab: 'Tab',
    form: 'Form',
    createArticle: 'Create Article',
    editArticle: 'Edit Article',
    articleList: 'Article List',
    errorPages: 'Error Pages',
    page401: '401',
    page404: '404',
    errorLog: 'Error Log',
    excel: 'Excel',
    exportExcel: 'Export Excel',
    selectExcel: 'Export Selected',
    mergeHeader: 'Merge Header',
    uploadExcel: 'Upload Excel',
    zip: 'Zip',
    pdf: 'PDF',
    exportZip: 'Export Zip',
    theme: 'Theme',
    clipboardDemo: 'Clipboard',
    i18n: 'I18n',
    externalLink: 'External Link',
    profile: 'Profile',
    pipelineTable: 'Pipeline Table',
    pipelineMonitor: 'Pipeline Monitor'
  },
  navbar: {
    dashboard: 'Dashboard',
    github: 'Github',
    logOut: 'Log Out',
    profile: 'Profile',
    theme: 'Theme',
    size: 'Global Size'
  },
  login: {
    title: 'Login Form',
    logIn: 'Login',
    username: 'Username',
    password: 'Password',
    any: 'any',
    thirdparty: 'Or connect with',
    thirdpartyTips: 'Can not be simulated on local, so please combine you own business simulation! ! !'
  },
  documentation: {
    documentation: 'Documentation',
    github: 'Github Repository'
  },
  permission: {
    addRole: 'New Role',
    editPermission: 'Edit',
    roles: 'Your roles',
    switchRoles: 'Switch roles',
    tips: 'In some cases, using v-permission will have no effect. For example: Element-UI  el-tab or el-table-column and other scenes that dynamically render dom. You can only do this with v-if.',
    delete: 'Delete',
    confirm: 'Confirm',
    cancel: 'Cancel'
  },
  guide: {
    description: 'The guide page is useful for some people who entered the project for the first time. You can briefly introduce the features of the project. Demo is based on ',
    button: 'Show Guide'
  },
  components: {
    documentation: 'Documentation',
    tinymceTips: 'Rich text is a core feature of the management backend, but at the same time it is a place with lots of pits. In the process of selecting rich texts, I also took a lot of detours. The common rich texts on the market have been basically used, and I finally chose Tinymce. See the more detailed rich text comparison and introduction.',
    dropzoneTips: 'Because my business has special needs, and has to upload images to qiniu, so instead of a third party, I chose encapsulate it by myself. It is very simple, you can see the detail code in @/components/Dropzone.',
    stickyTips: 'when the page is scrolled to the preset position will be sticky on the top.',
    backToTopTips1: 'When the page is scrolled to the specified position, the Back to Top button appears in the lower right corner',
    backToTopTips2: 'You can customize the style of the button, show / hide, height of appearance, height of the return. If you need a text prompt, you can use element-ui el-tooltip elements externally',
    imageUploadTips: 'Since I was using only the vue@1 version, and it is not compatible with mockjs at the moment, I modified it myself, and if you are going to use it, it is better to use official version.'
  },
  table: {
    dynamicTips1: 'Fixed header, sorted by header order',
    dynamicTips2: 'Not fixed header, sorted by click order',
    dragTips1: 'The default order',
    dragTips2: 'The after dragging order',
    title: 'Title',
    importance: 'Imp',
    type: 'Type',
    remark: 'Remark',
    search: 'Search',
    add: 'Add',
    export: 'Export',
    reviewer: 'reviewer',
    id: 'ID',
    date: 'Date',
    author: 'Author',
    readings: 'Readings',
    status: 'Status',
    actions: 'Actions',
    edit: 'Edit',
    publish: 'Publish',
    draft: 'Draft',
    delete: 'Delete',
    cancel: 'Cancel',
    confirm: 'Confirm'
  },
  example: {
    warning: 'Creating and editing pages cannot be cached by keep-alive because keep-alive include does not currently support caching based on routes, so it is currently cached based on component name. If you want to achieve a similar caching effect, you can use a browser caching scheme such as localStorage. Or do not use keep-alive include to cache all pages directly. See details'
  },
  errorLog: {
    tips: 'Please click the bug icon in the upper right corner',
    description: 'Now the management system are basically the form of the spa, it enhances the user experience, but it also increases the possibility of page problems, a small negligence may lead to the entire page deadlock. Fortunately Vue provides a way to catch handling exceptions, where you can handle errors or report exceptions.',
    documentation: 'Document introduction'
  },
  excel: {
    export: 'Export',
    selectedExport: 'Export Selected Items',
    placeholder: 'Please enter the file name (default excel-list)'
  },
  zip: {
    export: 'Export',
    placeholder: 'Please enter the file name (default file)'
  },
  pdf: {
    tips: 'Here we use window.print() to implement the feature of downloading PDF.'
  },
  theme: {
    change: 'Change Theme',
    documentation: 'Theme documentation',
    tips: 'Tips: It is different from the theme-pick on the navbar is two different skinning methods, each with different application scenarios. Refer to the documentation for details.'
  },
  tagsView: {
    refresh: 'Refresh',
    close: 'Close',
    closeOthers: 'Close Others',
    closeAll: 'Close All'
  },
  settings: {
    title: 'Page style setting',
    theme: 'Theme Color',
    tagsView: 'Open Tags-View',
    fixedHeader: 'Fixed Header',
    sidebarLogo: 'Sidebar Logo'
  },
  pipeline: {
    name: 'Pipeline Name',
    aliasName: 'Alias Name',
    status: 'Status',
    statusValues: {
      started: 'Started',
      starting: 'Starting',
      stopped: 'Stopped',
      stopping: 'Stopping'
    },
    statusOptions: [
      {
        key: 'run',
        value: 'RUN'
      },
      {
        key: 'stop',
        value: 'STOP'
      }
    ],
    infoStatusMap: {
      run: {
        key: 'run',
        value: 'RUN',
        type: 'success'
      },
      stop: {
        key: 'stop',
        value: 'STOP',
        type: 'info'
      },
      stopping: {
        key: 'stopping',
        value: 'STOPPING',
        type: 'warning'
      },
      scheduling: {
        key: 'scheduling',
        value: 'SCHEDULING',
        type: 'warning'
      },
      scheduled: {
        key: 'scheduled',
        value: 'SCHEDULED',
        type: ''
      },
      deleting: {
        key: 'deleting',
        value: 'DELETING',
        type: 'danger'
      }
    },
    createTime: 'Create Time',
    mysqlAddress: 'Mysql Address',
    mysqlPort: 'Mysql Port',
    mysqlUser: 'Mysql User',
    mysqlPassword: 'Mysql Password',
    mysqlServerId: 'Mysql Slave Server ID',
    mysqlFlavor: 'Flavor',
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
    mysqlMode: 'Mysql Binlog Mode',
    mysqlModeOptions: [
      {
        key: 'position',
        value: 'Common'
      },
      {
        key: 'gtid',
        value: 'GTID'
      }
    ],
    output: {
      sender: {
        name: 'Sender Name',
        type: 'Sender Type',
        typeOptions: [
          {
            key: 'stdout',
            value: 'Stdout'
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
          },
          {
            key: 'elastic',
            value: 'ElasticSearch'
          }
        ]
      }
    },
    bindNode: 'Bind Node',
    remark: 'Remark',
    filter: {
      type: 'Type',
      rule: 'Rule'
    },
    filterValues: {
      white: 'White List',
      black: 'Black List'
    }
  },
  node: {
    createTime: 'Join Time',
    name: 'Node Name',
    leader: 'Leader Node',
    ip: 'Current IP',
    status: 'Node Status',
    version: 'Node Version',
    roleMap: {
      leader: 'Leader',
      follower: 'Follower'
    }
  },
  capacity: {
    cpuCores: 'CPU Cores',
    memory: 'Memory',
    cpuUsage: 'CPU Usage',
    memoryUsage: 'Memory Usage'
  },
  pipeline_table: {
    select: {
      pleaseSelect: 'Please select'
    },
    date: {
      pickDate: 'Please select'
    },
    input: {
      pleaseInput: 'Please input'
    },
    filter: {
      name: 'Filter',
      delete: 'Delete',
      addFilter: 'New Filter',
      whiteList: 'White List',
      blackList: 'Black List',
      place: 'Enter in this format: database.table or database'
    },
    kafka: {
      brokers: 'Multiple brokers are separated by commas.'
    },
    pipeline: {
      name: 'Pipeline name, globally unique',
      aliasName: 'Pipeline alias, optional, easy to understand',
      run: 'Run'
    },
    http: {
      retries: 'Retries'
    },
    rabbit: {
      tips: 'Configure rabbitmq exchange, and binlogo will send a message to the exchange. The routing key is the database plus data table, such as database.table. Users need to create a corresponding rabbitmq queue when using it'
    },
    redis: {
      tips: 'Binlogo will send data to the List of Redis. The command used is RPUSH -- Insert data from the right side of the list'
    },
    rocket: {
      tips: 'Because user generally uses Alibaba cloud services directly, rocketMQ of Alibaba cloud configurations and reference addresses are provided here'
    },
    tips: {
      fix_pos_newest: 'If this is true, binlogo will use the newest positon when could not find binlog file index'
    },
    actions: 'Actions',
    status: 'Status',
    edit: 'Edit',
    delete: 'Delete',
    run: 'Run it',
    stop: 'Stop',
    name: 'Pipeline name',
    server_id: 'If it is blank, the background will be generated automatically',
    // sortOptions: [{ label: 'Time Ascending', key: '+id' }, { label: 'Time Descending', key: '-id' }]
    sortOptions: [{ label: 'Ascending', key: '+id' }, { label: 'Descending', key: '-id' }],
    statusOptions: [{ label: 'Run', key: 'run' }, { label: 'Stop', key: 'stop' }]
  },
  pipe_detail: {
    pipeline: 'Pipeline',
    remark: 'Remark',
    detail: 'Detail',
    runStatus: 'Run Status',
    instance: 'Instance',
    govern: 'Govern',
    filter: 'Filter',
    event: 'Event',
    filterTip: 'Enter the table name to verify whether the database or data table can pass under the current filtering rule. The table name format is database name plus table name, such as test_ db.test_tbl',
    validRule: 'Database name or data table name to be verified',
    button_valid: 'Valid',
    validTipTrue: ' will be filtered',
    validTipFalse: ' will not be filtered',
    positionTip: 'Manually select a binlog synchronization location. For fault recovery or special business requirements.',
    eventTip: 'Displays the last 20 event records',
    event_table: {
      first_time: 'First Time',
      last_time: 'Last Time',
      type: 'Type',
      count: 'AGG Times',
      node_name: 'From Node',
      node_ip: 'Node IP',
      message: 'Content'
    }
  },
  node_table: {
    statusMap: {
      ready: {
        yes: 'Ready',
        no: 'Not Ready'
      },
      network_unavailable: {
        yes: 'Network Unavailable'
      }
    },
    name: 'Node Name',
    ready: 'Status',
    sortOptions: [{ label: 'Ascending', key: '+id' }, { label: 'Descending', key: '-id' }],
    readyOptions: [{ label: 'Ready', key: 'yes' }, { label: 'Not Ready', key: 'no' }]
  },
  instance: {
    create_time: 'Start Time',
    pipeline_name: 'Pipeline',
    node_name: 'Node Name'
  },
  register: {
    create_time: 'Registration Time',
    name: 'Node Name',
    ip: 'Node IP',
    version: 'Version'
  },
  cluster: {
    title: 'Cluster Info',
    name: 'Cluster Name',
    tabMap: {
      election: {
        key: 'election',
        label: 'Election Node'
      },
      register: {
        key: 'register',
        label: 'Registration Node'
      },
      instance: {
        key: 'instance',
        label: 'Instance '
      }
    },
    instanceTip: 'Displays the currently running pipeline instance.',
    registerTip: 'Displays all currently registered (online) nodes.',
    electionTip: 'Displays the status of the master node.'
  },
  global: {
    submit: 'Submit'
  }
}
