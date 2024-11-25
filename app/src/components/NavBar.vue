<script lang="ts" setup>

import Menubar from 'primevue/menubar'
import { ref } from 'vue';


const items = ref([
    {
        label: 'Продукты',
        icon: 'pi pi-box',
        route: '/products'
    }
]);

</script>

<template>
    <Menubar :model="items">
        <template #item="{ item, props, hasSubmenu }">
            <router-link v-if="item.route" v-slot="{ href, navigate }" :to="item.route" custom>
                <a v-ripple :href="href" v-bind="props.action" @click="navigate">
                    <span :class="item.icon" />
                    <span>{{ item.label }}</span>
                </a>
            </router-link>
            <a v-else v-ripple :href="item.url" :target="item.target" v-bind="props.action">
                <span :class="item.icon" />
                <span>{{ item.label }}</span>
                <span v-if="hasSubmenu" class="pi pi-fw pi-angle-down" />
            </a>
        </template>
    </Menubar>
</template>


<style></style>