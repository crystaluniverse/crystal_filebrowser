<template>
    <div id="document">
</div>
</template>

<script>
import { mapState } from 'vuex'
import url from '@/utils/url'
import { baseURL } from '@/utils/constants'
import { files as api } from '@/api'
import crypto from 'crypto'

export default {
  name: 'document',
  components: {
  },
  props: ['fileData', 'user'],
  data: function () {
    return {
    }
  },
  computed: {
    ...mapState(['jwt']),
    documentLocation () {
      return `${baseURL}/api/raw${this.fileData.path}?auth=${this.jwt}`
    },
  },
  mounted () {
    this.renderDocument()
  },
  methods: {
    renderDocument () {
      const config = {
            "document": {
                "fileType": this.fileData.extension.replace('.', ''),
                "key": this.fileData.key.replace(new RegExp('[^A-Za-z0-9]+', 'g'), '').slice(0, 19),
                "title": this.fileData.name,
                "url": `${window.location.origin}${this.documentLocation}`
            },
            "editorConfig": {
              "callbackUrl": `${window.location.origin}/callback?auth=${this.jwt}&filepath=${this.fileData.path}`,
              "user": {
                "id": crypto.randomBytes(8).toString('hex'),
                "name": this.user.username
              }
            }
        }
        new window.DocsAPI.DocEditor("document", config)
        console.log(config)
    }
  }
}
</script>
