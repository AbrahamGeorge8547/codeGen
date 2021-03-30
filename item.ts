import * as mongoose from 'mongoose';


const ItemSchema = new mongoose.Schema({
referenceUrl: String,
website: String,
email: String,
price: Number,
prices: Array,
})

export interface IItem extends mongoose.Document {
referenceUrl: string,
website: string,
email: string,
price: number,
prices: number[],
}

export default mongoose.model<IItem>("Item", ItemSchema)