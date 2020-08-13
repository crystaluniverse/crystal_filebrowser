<template>
  <div class="card floating" id="share">
    <div class="card-title">
      <h2>Share</h2>
    </div>
    <div class="card-content">
      <v-data-table
        dark
        :headers="userheaders"
        :items="userpermissions"
        :search="userSearch"
        item-key="name"
        show-select
        v-model="selectedUsers"
        :footer-props="{ 'items-per-page-options': [5, 10, 15, -1] }"
        :items-per-page="5"
      >
        <template v-slot:top>
          <v-toolbar flat>
            <v-row>
              <v-col cols="4">
                People who have access
              </v-col>
              <v-spacer></v-spacer>
              <v-col cols="5">
                <v-text-field
                  v-model="userSearch"
                  append-icon="mdi-magnify"
                  label="Search"
                  single-line
                  hide-details
                ></v-text-field>
                <v-btn class="ml-2" :disabled="selectedUsersCount == 0"
                  ><i class="material-icons" dark small>delete</i>
                  <span v-if="selectedUsersCount > 0">{{
                    selectedUsersCount
                  }}</span></v-btn
                >
              </v-col>
            </v-row>
          </v-toolbar>
        </template>
        <template v-slot:item.permission="{ item }">
          {{ permissionToHumanReadable(item.permission)}}
        </template>
        <template v-slot:item.action="{ item }">
          <i
            class="material-icons"
            dark
            small
            @click="deleteUserAccess(item.name)"
            >delete</i
          >
        </template>
      </v-data-table>

      <v-row dense align="center">
        <v-col cols="8">
          <v-text-field
            dark
            v-model="users"
            clearable
            label="Insert 3bot names comma seperated"
          >
          </v-text-field>
        </v-col>
        <v-col cols="2">
          <v-select
            class="permissions"
            dark
            v-model="sharePermission"
            :items="permissionList"
          >
          </v-select>
        </v-col>
        <v-col cols="2" align="center">
          <button class="button button--flat" @click="shareUsers">
            Add
          </button>
        </v-col>
      </v-row>

      <v-data-table
        dark
        :headers="linkheaders"
        :items="existingLinks"
        item-key="uuid"
        show-select
        v-model="selectedLinks"
        :footer-props="{ 'items-per-page-options': [5, 10, 15, -1] }"
        :items-per-page="5"
      >
        <template v-slot:top>
          <v-toolbar flat>
            <v-row align="center">
              <v-col cols="4">
                People who have access
              </v-col>
              <v-spacer></v-spacer>
              <v-btn class="mr-3" :disabled="selectedLinksCount == 0"
                ><i class="material-icons" dark small>delete</i>
                <span v-if="selectedLinksCount > 0">{{
                  selectedLinksCount
                }}</span></v-btn
              >
            </v-row>
          </v-toolbar>
        </template>
        <template v-slot:item.permission="{ item }">
          {{ permissionToHumanReadable(item.permission)}}
        </template>
        <template v-slot:item.action="{ item }">
          <i class="material-icons" dark small @click="deleteLink(item.uuid)"
            >delete</i
          >
        </template>
      </v-data-table>
      <div class="mt-2">Create a link with the following rights:</div>
      <v-row align="center" dense>
        <v-col cols="2">
          <v-select
            dense
            class="permissions"
            dark
            v-model="sharePermission"
            :items="permissionList"
          >
          </v-select>
        </v-col>
        <v-col cols="2">
          <button class="buteton button--flat" @click="getLink">
            Get link
          </button>
        </v-col>
        <v-spacer></v-spacer>
      </v-row>
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

      userheaders: [
        {
          text: "3bot name",
          value: "name",
        },
        {
          text: "Acces right",
          value: "permission",
          filterable: true,
        },
        {
          text: "Action",
          value: "action",
        },
      ],
      users: "",
      sharePermission: "",
      userpermissions: [
        { name: "hamdy", permission: "rw-" },
        { name: "tobias", permission: "r--" },
        { name: "ivan", permission: "r--" },
        { name: "alex", permission: "rwd" },
        { name: "mathias", permission: "r--" },
        { name: "jonas", permission: "r--" },
      ],
      selectedUsers: [],
      userSearch: "",

      linkheaders: [
        { text: "URL", value: "link" },
        { text: "Acces right", value: "permission" },
        { text: "Action", value: "action" },
      ],
      existingLinks: [
        {
          link: "https://threefold.io/DDN4L34M6N434RJ3",
          permission: "rw-",
          uuid: "blah",
        },
        {
          link: "https://threefold.io/DDN4L34M6N434RJ3",
          permission: "rw-",
          uuid: "ads",
        },
        {
          link: "https://threefold.io/DDN4L34M6N434RJ3",
          permission: "rw-",
          uuid: "fdf",
        },
        {
          link: "https://threefold.io/DDN4L34M6N434RJ3",
          permission: "rw-",
          uuid: "qw",
        },
        {
          link: "https://threefold.io/DDN4L34M6N434RJ3",
          permission: "rw-",
          uuid: "AE",
        },
        {
          link: "https://threefold.io/DDN4L34M6N434RJ3",
          permission: "rw-",
          uuid: "ff",
        },
      ],
      linkPermission: "",
      selectedLinks: [],
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
    permissionList() {
      return Array.from(this.sharePermissions, (x) => x.name);
    },
    selectedUsersCount() {
      return this.selectedUsers.length;
    },
    selectedLinksCount() {
      return this.selectedLinks.length;
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
    this.sharePermission = this.sharePermissions[0].name;
    this.linkPermission = this.sharePermissions[0].name;

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
    mapUserPermission(users, permission){
      return users.map(user => {
        return {
          "name": user,
          "permission": permission
          }
      })
    },
    shareUsers: async function() {
      let users = this.users.split(",");
      users = this.mapUserPermission(users, this.sharePermission)
      console.log(users)
      try {
        await api.shareWithUsers(this.url, users);
      } catch (e) {
        console.log(e);
        this.$showError(e);
      }
    },
    deleteUserAccess: async function(user) {
      user = this.mapUserPermission([user], "")
      console.log(user)
      await api.shareWithUsers(this.url, user );
    },
    deleteAllShares: async function() {
      await api.deleteAllShares(this.url);
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
