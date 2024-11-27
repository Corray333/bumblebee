import { useToast } from "primevue"
import { ProductTransport } from "./transport"

export class ProductService{

    static async getProducts(){
        try {
            return ProductTransport.getProducts()
        } catch (error) {
            console.log("Err: ",error)
            useToast().add({severity:'error', summary: 'Не получилось получить продукты', life: 3000})
            return []
        }
    }
}