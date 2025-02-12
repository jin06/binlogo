import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'/'el-icon-x' the icon show in the sidebar
    noCache: true                if set true, the page will no be cached(default is false)
    affix: true                  if set true, the tag will affix in the tags-view
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  {
    path: '/redirect',
    component: Layout,
    hidden: true,
    children: [
      {
        path: '/redirect/:path(.*)',
        component: () => import('@/views/redirect/index')
      }
    ]
  },
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/auth-redirect',
    component: () => import('@/views/login/auth-redirect'),
    hidden: true
  },
  {
    path: '/404',
    component: () => import('@/views/error-page/404'),
    hidden: true
  },
  {
    path: '/401',
    component: () => import('@/views/error-page/401'),
    hidden: true
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        component: () => import('@/views/cluster/index'),
        name: 'Dashboard',
        meta: { title: 'dashboard', icon: 'el-icon-s-home', affix: true } // 使用 el-icon-s-home 图标
      }
    ]
  },
  {
    path: '/pipeline',
    component: Layout,
    redirect: '/pipeline/pipeline-table',
    name: 'pipeline',
    meta: {
      title: 'pipeline',
      icon: 'el-icon-s-operation' // 使用 el-icon-s-operation 图标
    },
    children: [
      {
        path: 'pipeline-table',
        component: () => import('@/views/pipeline/pipeline-table'),
        name: 'pipelineTable',
        meta: { title: 'pipelineTable', icon: 'el-icon-s-marketing' } // 使用 el-icon-s-marketing 图标
      },
      {
        path: 'pipeline-detail/:name',
        component: () => import('@/views/pipeline/detail/index'),
        name: 'pipelineDetail',
        meta: { title: 'pipelineDetail' },
        hidden: true
      }
    ]
  },
  {
    path: '/node',
    component: Layout,
    redirect: '/node/node-table',
    name: 'node',
    meta: {
      title: 'node',
      icon: 'el-icon-s-platform' // 使用 el-icon-s-platform 图标
    },
    children: [
      {
        path: 'node-table',
        component: () => import('@/views/node/node-table'),
        name: 'nodeTable',
        meta: { title: 'nodeTable', icon: 'el-icon-s-data' } // 使用 el-icon-s-data 图标
      }
    ]
  },
  {
    path: '/cluster',
    component: Layout,
    redirect: '/cluster/index',
    name: 'cluster',
    meta: {
      title: 'cluster',
      icon: 'el-icon-s-tools' // 使用 el-icon-s-tools 图标
    },
    children: [
      {
        path: 'cluster',
        component: () => import('@/views/cluster/index'),
        name: 'cluster',
        meta: { title: 'cluster', icon: 'el-icon-s-management' } // 使用 el-icon-s-management 图标
      }
    ]
  }
]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user roles
 */
export const asyncRoutes = [

  /** when your routing map is too long, you can split it into small modules **/

  // 404 page must be placed at the end !!!
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
