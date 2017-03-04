<template>
  <el-row id="app" :gutter="20">
    <el-col :span="6">
      <div class="grid-content bg-purple">
        <h3>Nav Pop Over</h3>
        <p>Pop over nav side column</p>
      </div>
    </el-col>

    <el-col :span="18">
      <div class="grid-content bg-purple">
        <h1>Main Content Here</h1>
        <router-view></router-view>
      </div>
    </el-col>
  </el-row>
</template>

<script>
 export default {
   name: 'app',
 };
</script>

<style>
 /* #app { */
 /* font-family: 'Avenir', Helvetica, Arial, sans-serif; */
 /* -webkit-font-smoothing: antialiased; */
 /* -moz-osx-font-smoothing: grayscale; */
 /* text-align: center; */
 /* color: #2c3e50; */
 /* margin-top: 60px; */
 /* } */
</style>
