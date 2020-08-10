<template>
  <div class="card floating" id="share">
    <div class="card-title">
      <h2>Share</h2>
    </div>

    <div class="card-content">
      <!--<li v-if="!hasPermanent">
          <a @click="getPermalink" :aria-label="$t('buttons.permalink')">{{ $t('buttons.permalink') }}</a>
        </li> -->
      <h3>Share with people</h3>
      <input
        id="nameInput"
        type="text"
        v-focus
        v-model="users"
        placeholder="3bot names (comma separated)"
      />
      <select class="permissions" v-model="sharePermission">
        <option
          v-for="option in sharePermissions"
          :key="option.value"
          :value="option.value"
          >{{ option.name }}</option
        >
      </select>

      <button class="button button--flat" @click="shareUsers">
        Submit
      </button>
      <h4>Already shared with</h4>
      <div v-for="permission in sharePermissions" :key="permission.value">
        <h5>
          {{ permissionToHumanReadable(permission.value) }}
        </h5>
        <div v-for="user in userpermissions[permission.value]" :key="user">
          <i class="material-icons" @click="deleteUserAccess(user, permission.value)">
            delete
          </i>
          {{ user }}
        </div>
      </div>
      <h3>Get shareable link</h3>
      <select class="permissions" v-model="sharePermission">
        <option
          v-for="option in sharePermissions"
          :key="option.value"
          :value="option.value"
          >{{ option.name }}</option
        >
      </select>
      <button class="button button--flat" @click="getLink">
        Get link
      </button>

      <h4>Existing links</h4>
      <div v-for="link in existingLinks" :key="link.uuid">
        <i class="material-icons" @click="deleteLink(link.uuid)">
          delete
        </i>
        {{ link.link }} {{ permissionToHumanReadable(link.permission) }}
      </div>
    </div>

    <div class="card-action">
      <button
        class="button button--flat"
        @click="$store.commit('closeHovers')"
        :aria-label="$t('buttons.close')"
        :title="$t('buttons.close')"
      >
        {{ $t("buttons.close") }}
      </button>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import { share as api } from "@/api";
import { baseURL } from "@/utils/constants";
import moment from "moment";
import Clipboard from "clipboard";

export default {
  name: "share",
  data: function() {
    return {
      time: "",
      unit: "hours",
      hasPermanent: false,
      links: [],
      clip: null,
      sharePermissions: [
        {
          name: "Viewer",
          value: "r--",
        },
        {
          name: "Editor",
          value: "rw-",
        },
        {
          name: "Owner",
          value: "rwd",
        },
      ],
      users: "",
      sharePermission: "",
      linkPermission: "",
      userpermissions: { rwd: ["hamdy", "tobias"], "r--": ["aaa"] },
      existingLinks: [{ link: "url", permission: "rw-", uuid: "blah" }],
    };
  },
  computed: {
    ...mapState(["req", "selected", "selectedCount"]),
    ...mapGetters(["isListing"]),
    url() {
      if (!this.isListing) {
        return this.$route.path;
      }

      if (this.selectedCount === 0 || this.selectedCount > 1) {
        // This shouldn't happen.
        return;
      }

      return this.req.items[this.selected[0]].url;
    },
  },
  async beforeMount() {
    try {
      const links = await api.get(this.url);
      this.links = links;
      this.sort();

      for (let link of this.links) {
        if (link.expire === 0) {
          this.hasPermanent = true;
          break;
        }
      }
    } catch (e) {
      this.$showError(e);
    }
  },
  mounted() {
    this.clip = new Clipboard(".copy-clipboard");
    this.clip.on("success", () => {
      this.$showSuccess(this.$t("success.linkCopied"));
    });
    this.sharePermission = this.sharePermissions[0].value;
    this.linkPermission = this.sharePermissions[0].value;

    // this.userpermissions = api.listUserpermissions(this.url)
    // this.existingLinks = api.listLinks(this.url)

  },
  beforeDestroy() {
    this.clip.destroy();
  },
  methods: {
    permissionToHumanReadable(permission) {
      return this.sharePermissions.find((perm) => permission == perm.value)
        .name;
    },
    submit: async function() {
      if (!this.time) return;

      try {
        const res = await api.create(this.url, this.time, this.unit);
        this.links.push(res);
        this.sort();
      } catch (e) {
        this.$showError(e);
      }
    },
    getPermalink: async function() {
      try {
        const res = await api.create(this.url);
        this.links.push(res);
        this.sort();
        this.hasPermanent = true;
      } catch (e) {
        this.$showError(e);
      }
    },
    // deleteLink: async function(event, link) {
    //   event.preventDefault();
    //   try {
    //     await api.remove(link.hash);
    //     if (link.expire === 0) this.hasPermanent = false;
    //     this.links = this.links.filter((item) => item.hash !== link.hash);
    //   } catch (e) {
    //     this.$showError(e);
    //   }
    // },
    shareUsers: async function() {
      const users = this.users.split(",");
      try {
        await api.shareWithUsers(this.url, users, this.sharePermission);
      } catch (e) {
        console.log(e);
        this.$showError(e);
      }
    },
    deleteUserAccess: async function(user, permission){
      await api.deleteUserAccess(this.url, user,permission)
    },
    deleteAllShares: async function(){
      await api.deleteAllShares(this.url)
    },
    getLink: async function() {
      await api.getShareableLink(this.url, this.linkPermission);
    },
    deleteLink: async function(uuid) {
      await api.deleteSharableLink(uuid);
    },
    humanTime(time) {
      return moment(time * 1000).fromNow();
    },
    buildLink(hash) {
      return `${window.location.origin}${baseURL}/share/${hash}`;
    },
    sort() {
      this.links = this.links.sort((a, b) => {
        if (a.expire === 0) return -1;
        if (b.expire === 0) return 1;
        return new Date(a.expire) - new Date(b.expire);
      });
    },
  },
};
</script>
