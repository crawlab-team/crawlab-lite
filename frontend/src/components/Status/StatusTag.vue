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
          PROCESSING: { label: 'Processing', type: 'primary' },
          RUNNING: { label: 'Running', type: 'warning' },
          FINISHED: { label: 'Finished', type: 'success' },
          ERROR: { label: 'Error', type: 'danger' },
          CANCELLED: { label: 'Cancelled', type: 'info' },
          ABNORMAL: { label: 'Abnormal', type: 'danger' }
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
        if (this.status === 'finished') {
          return 'el-icon-check'
        } else if (this.status === 'pending') {
          return 'el-icon-loading'
        } else if (this.status === 'running') {
          return 'el-icon-loading'
        } else if (this.status === 'error') {
          return 'el-icon-error'
        } else if (this.status === 'cancelled') {
          return 'el-icon-video-pause'
        } else if (this.status === 'abnormal') {
          return 'el-icon-warning'
        } else {
          return 'el-icon-question'
        }
      }
    }
  }
</script>

<style scoped>

</style>
