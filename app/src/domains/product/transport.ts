import {api} from '@/api'
import type { Product } from './entities'

export class ProductTransport{
    static async getProducts(): Promise<Product[]>{
        try {
            const response = await api.get('/products')
            console.log(response.data)
            return response.data
        } catch (error) {
            console.log(error)
            return []
        }
    }

    static async createProduct(product: Product, photo: File): Promise<Product>{
        const formData = new FormData();
        formData.append('photo', photo);
        formData.append('product', JSON.stringify(product));
        const response = await api.post('/products', formData)
        return response.data
    }
    static async updateProduct(product: Product, photo: File): Promise<Product>{
        const formData = new FormData();
        formData.append('photo', photo);
        formData.append('product', JSON.stringify(product));
        const response = await api.put(`/products/${product.id}`, formData)
        return response.data
    }
    static async deleteProduct(id: number): Promise<void>{
        await api.delete(`/products/${id}`)
    }

    static async reorderProduct(id: number, position: number): Promise<void>{
        await api.put(`/products/${id}/reorder?new_position=${position}`)
    }
}