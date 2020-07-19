<template>
  <el-tag :type="type" class="status-tag">
    <i :class="icon" />
    {{ $t(label) }}
  </el-tag>
</template>

<script>
  export default {
    name: 'StatusTag',
    props: {
      status: {
        type: String,
        default: ''
      }
    },
    data() {
      return {
        statusDict: {
          PENDING: { label: 'Pending', type: 'primary' },
          RUNNING: { label: 'Running', type: 'warning' },
          FINISHED: { label: 'Finished', type: 'success' },
          ERROR: { label: 'Error', type: 'danger' },
          CANCELLED: { label: 'Cancelled', type: 'info' }
        }
      }
    },
    computed: {
      type() {
        const s = this.statusDict[this.status]
        if (s) {
          return s.type
        }
        return ''
      },
      label() {
        const s = this.statusDict[this.status]
        if (s) {
          return s.label
        }
        return 'NA'
      },
      icon() {
        if (this.status === 'FINISHED') {
          return 'el-icon-check'
        } else if (this.status === 'PENDING') {
          return 'el-icon-loading'
        } else if (this.status === 'RUNNING') {
          return 'el-icon-loading'
        } else if (this.status === 'ERROR') {
          return 'el-icon-error'
        } else if (this.status === 'CANCELLED') {
          return 'el-icon-video-pause'
        } else {
          return 'el-icon-question'
        }
      }
    }
  }
</script>

<style scoped>

</style>
