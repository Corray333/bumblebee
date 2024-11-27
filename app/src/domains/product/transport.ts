import {api} from '@/api'
import type { Product } from './entities'
import { useToast } from 'primevue'
import { useSystemStore } from '@/stores/counter'

const systemStore = useSystemStore()

export class ProductTransport{
    static async getProducts(): Promise<Product[]>{
        try {
            const response = await api.get('/products')
            return response.data
        } catch (error) {
            console.log("Err: ",error)
            if (systemStore.toast) systemStore.toast.add({severity:'error', summary: 'Не получилось получить продукты', life: 3000})
            return []
        }
    }

    static async createProduct(product: Product, photo: File | undefined){
        try {     
            const formData = new FormData();
            if (!photo) throw new Error('Фото не выбрано');
            formData.append('photo', photo);
            formData.append('product', JSON.stringify(product));
            await api.post('/products', formData)
            return 
        } catch (error) {
            console.log(error)
            if (systemStore.toast) systemStore.toast.add({severity:'error', summary: 'Не получилось создать продукт', life: 3000})
        }
    }
    
    static async updateProduct(product: Product, photo: File | undefined){
        try {
            const formData = new FormData();
            if(photo)formData.append('photo', photo);
            formData.append('product', JSON.stringify(product));
            await api.put(`/products/${product.id}`, formData)
        } catch (error) {
            console.log(error)
            if (systemStore.toast) systemStore.toast.add({severity:'error', summary: 'Не получилось обновить продукт', life: 3000})
        }

    }
    static async deleteProduct(id: number): Promise<void>{
        try {
            await api.delete(`/products/${id}`)
        } catch (error) {
            console.log(error)
            if (systemStore.toast) systemStore.toast.add({severity:'error', summary: 'Не получилось удалить продукт', life: 3000})
        }
    }

    static async reorderProduct(id: number, position: number): Promise<void>{
        try {
            await api.put(`/products/${id}/reorder?new_position=${position}`)
        } catch (error) {
            console.log(error)
            if (systemStore.toast) systemStore.toast.add({severity:'error', summary: 'Не получилось изменить позицию продукта', life: 3000})
        }
    }
}