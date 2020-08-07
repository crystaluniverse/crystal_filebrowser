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
      <!-- <li v-for="link in links" :key="link.hash">
          <button class="action"
            @click="deleteLink($event, link)"
            :aria-label="$t('buttons.delete')"
            :title="$t('buttons.delete')"><i class="material-icons">delete</i></button>

          <button class="action copy-clipboard"
            :data-clipboard-text="buildLink(link.hash)"
            :aria-label="$t('buttons.copyToClipboard')"
            :title="$t('buttons.copyToClipboard')"><i class="material-icons">content_paste</i></button>
        </li>

        <li>
          <input 
            type="number"
            max="2147483647"
            min="0"
            @keyup.enter="submit"
            v-model.trim="time">
          <select v-model="unit" :aria-label="$t('time.unit')">
            <option value="seconds">{{ $t('time.seconds') }}</option>
            <option value="minutes">{{ $t('time.minutes') }}</option>
            <option value="hours">{{ $t('time.hours') }}</option>
            <option value="days">{{ $t('time.days') }}</option>
          </select>
          <button class="action"
            @click="submit"
            :aria-label="$t('buttons.create')"
            :title="$t('buttons.create')"><i class="material-icons">add</i></button>
        </li>
      </ul>-->
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
      linkPermission: ""
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
    this.sharePermission = this.sharePermissions[0].value
    this.linkPermission = this.sharePermissions[0].value
  },
  beforeDestroy() {
    this.clip.destroy();
  },
  methods: {
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
    deleteLink: async function(event, link) {
      event.preventDefault();
      try {
        await api.remove(link.hash);
        if (link.expire === 0) this.hasPermanent = false;
        this.links = this.links.filter((item) => item.hash !== link.hash);
      } catch (e) {
        this.$showError(e);
      }
    },
    shareUsers: async function() {
      const users = this.users.split(",");
      try {
        await api.shareWithUsers(this.url, users, this.sharePermission);
      } catch (e) {
        console.log(e);
        this.$showError(e);
      }
    },
    getLink: async function() {
      await api.getShareableLink(this.url, this.linkPermission)
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
