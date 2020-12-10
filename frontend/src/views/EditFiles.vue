<template>
  <documentEditor v-if="request" :fileData="request" :user="user" />
</template>

<script>
import documentEditor from "@/components/files/DocumentEditor";
import { files as api } from '@/api'
import { mapGetters, mapState, mapMutations } from 'vuex'
function clean (path) {
  return path.endsWith('/') ? path.slice(0, -1) : path
}
export default {
  components: {
    documentEditor
  },
  mounted() {
    this.fetchData()
  },
  methods:{
     async fetchData () {
      let url = decodeURIComponent(this.$route.query.url)
      console.log(`url`, url)
      try {
        const res = await api.fetch(url)

        // if (clean(res.path) !== clean(`/${this.$route.params.pathMatch}`)) {
        //   return
        // }
        console.log(res)
        this.request = res
        // this.$store.commit('updateRequest', res)
        document.title = res.name
      } catch (e) {
        this.error = e
      } finally {
        console.log(`done`)
        // this.setLoading(false)
      }
    },
  },
  computed: {
    ...mapState([
      'req',
    ]),
    fileData() {
      let req = this.req
       return {
        path: req.path,
        extension: req.extension,
        key: req.key,
        name: req.name
      }
    },
    user() {
      // var url_string = window.location.href
      // var url = new URL(url_string);
      // var userName = url.searchParams.get("userName");

      let userName = "holbewoner"
      return {
        userName
      }
    },
  },
  data() {
    return {
      request: null
    };
  },
};
</script>
