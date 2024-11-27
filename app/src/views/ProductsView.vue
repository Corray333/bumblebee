<script setup lang="ts">

import { type Product } from '@/domains/product/entities';
import { ProductTransport } from '@/domains/product/transport';
import { Button, Column, DataTable, Dialog, FileUpload, Textarea } from 'primevue';
import { useToast } from 'primevue/usetoast';
import { onBeforeMount, ref } from 'vue';
const products = ref<Product[]>([])

onBeforeMount(async ()=>{
  products.value = await ProductTransport.getProducts()
})

const columns = ref([
    {field: 'id', header: 'ID'},
    {field: 'description', header: 'Описание'},
]);

const toast = useToast()

const onRowReorder = async (event: {dragIndex:number, dropIndex:number,value:Product[]}) => {
  const from = products.value[event.dragIndex].id
  products.value = event.value;
  await ProductTransport.reorderProduct(from, event.dropIndex+1)
  toast.add({severity:'success', summary: 'Порядок изменен', life: 3000});
};


const product = ref<Product>()
const productShowDialog = ref<boolean>(false)
const submitted = ref<boolean>(false)

const deleteProductsDialog = ref<boolean>(false)

const editProduct = (pick: Product) => {
    product.value = {...pick}
    productShowDialog.value = true
}

const hideDialog = () => {
    productShowDialog.value = false
    src.value = undefined
}

const saveProduct = async () => {
    submitted.value = true
    if (!product.value?.description) return
    if (product.value?.id == 0){
      await ProductTransport.createProduct(product.value, newFile.value)
      products.value = await ProductTransport.getProducts()
      toast.add({severity:'success', summary: 'Готово', detail: 'Продукт обновлен', life: 3000});
    } else {
      await ProductTransport.updateProduct(product.value, newFile.value)
      products.value = await ProductTransport.getProducts()
      toast.add({severity:'success', summary: 'Готово', detail: 'Продукт добавлен', life: 3000});
    }

    productShowDialog.value = false
    src.value = undefined
}

const newProduct = ()=>{
    product.value = {id: 0, description: '', img: ''}
    productShowDialog.value = true
}

const confirmDeleteProduct = (pick: Product) => {
    deleteProductsDialog.value = true
    product.value = pick
}

const deleteSelectedProducts = async () => {
  if (!product.value) return
    products.value = products.value.filter(val => val.id !== product.value?.id)
    deleteProductsDialog.value = false
    await ProductTransport.deleteProduct(product.value.id)
    product.value = undefined
    toast.add({severity:'success', summary: 'Готово', detail: 'Продукт удален', life: 3000});
}

const src = ref<string>()
const newFile = ref<File>()
function onFileSelect(event: { files: File[] }) {
  const file = event.files[0];
  newFile.value = file;
  const reader = new FileReader();

  reader.onload = async (e) => {
      if (e.target) src.value = e.target.result as string;
  };

  reader.readAsDataURL(file);
}


</script>

<template>


  <Dialog v-model:visible="deleteProductsDialog" :style="{ width: '550px' }" header="Удаление" :modal="true">
    <div class="flex items-center gap-4 w-full">
        <i class="pi pi-exclamation-triangle !text-3xl" />
        <span v-if="product" class=" text-ellipsis text-nowrap w-full overflow-hidden">Вы уверены, что хотите удалить продукт {{product.description}}</span>
    </div>

    <template #footer>
        <Button label="Отмена" icon="pi pi-times" text @click="deleteProductsDialog = false" />
        <Button label="Да" icon="pi pi-check" text @click="deleteSelectedProducts" />
    </template>
  </Dialog>

  <Dialog v-model:visible="productShowDialog" :style="{ width: '450px' }" header="Детали продукта" :modal="true">
    <div v-if="product" class="flex flex-col gap-6">
        <img v-if="!src && product.img" :src="product.img" :alt="product.img" class="shadow-md rounded-xl w-full" />

        <img v-if="src" :src="src" alt="Image" class="shadow-md rounded-xl w-full" />
        <div class="flex justify-center gap-2">
          <FileUpload mode="basic" choose-label="Выбрать" @select="onFileSelect" customUpload auto severity="secondary" class="p-button-outlined" />
          <Button v-if="src" @click="src = undefined" label="Отменить" icon="pi pi-times" />
        </div>

        <div>
            <label for="description" class="block font-bold mb-3">Описание</label>
            <Textarea id="description" v-model="product.description" required="true" rows="3" cols="20" class="w-full" fluid />
            <small v-if="submitted && !product?.description" class="text-red-500">Описание является обязательным параметром</small>
          </div>
    </div>

    <template #footer>
        <Button label="Отмена" icon="pi pi-times" text @click="hideDialog" />
        <Button label="Сохранить" icon="pi pi-check" @click="saveProduct" />
    </template>
  </Dialog>

  <main>
    <div class="flex justify-between items-center">
      <h1>Продукты</h1>
      <Button icon="pi pi-plus" label="Добавить" @click="newProduct()" />
    </div>
    <DataTable  :value="products" editMode="row" resizableColumns @rowReorder="onRowReorder" tableStyle="min-width: 50rem">
      <Column rowReorder headerStyle="width: 3rem" :reorderableColumn="false" />
      <Column v-for="col of columns" :field="col.field" :header="col.header" :key="col.field"></Column>
      <Column header="Картинка">
        <template #body="slotProps">
            <img :src="slotProps.data.img" :alt="slotProps.data.image" class="w-24 rounded" />
        </template>
      </Column>
      <Column :exportable="false" style="min-width: 12rem">
        <template #body="slotProps">
            <Button icon="pi pi-pencil" outlined rounded class="mr-2" @click="editProduct(slotProps.data)" />
            <Button icon="pi pi-trash" outlined rounded severity="danger" @click="confirmDeleteProduct(slotProps.data)" />
        </template>
      </Column>
    </DataTable>
  </main>
</template>

<style scoped>

main{
  @apply p-8
}

</style>