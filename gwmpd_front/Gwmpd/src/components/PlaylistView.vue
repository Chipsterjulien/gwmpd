<template lang="html">
  <div class="">
    PlaylistView
  </div>
</template>

<script>
export default {
  name: 'PlaylistView',
  mounted () {
    this.$resource('v1/allFiles').get().then((response) => {
      console.log(response.data)
    })
  }
}
</script>

<style lang="scss" scoped>
</style>
