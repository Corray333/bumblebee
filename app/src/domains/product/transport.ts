import {api} from '@/api'
import type { Product } from './entities'

export class ProductTransport{
    static async getProducts(): Promise<Product[]>{
        console.log(import.meta.env.VITE_API_URL)
        try {
            const response = await api.get('/products')
            console.log(response.data)
            return response.data
        } catch (error) {
            console.log(error)
            return []
        }
    }

    static async createProduct(product: Product): Promise<Product>{
        const response = await api.post('/products', product)
        return response.data
    }
    static async updateProduct(product: Product): Promise<Product>{
        const response = await api.put(`/products/${product.id}`, product)
        return response.data
    }
    static async deleteProduct(id: number): Promise<void>{
        await api.delete(`/products/${id}`)
    }

    static async reorderProduct(id: number, position: number): Promise<void>{
        await api.put(`/products/${id}/reorder?new_position=${position}`)
    }
}