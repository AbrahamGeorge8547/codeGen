import {Schema,   model} from 'mongoose';


const ItemSchema = new Schema({ 
	email: {required: true, type: String},  
	price: {required: true, type: Number},  
	prices: {required: true, type: Array},  
	referenceUrl: {required: false, type: Boolean},  
	website: {required: true, type: String} 
	})

export default model<IItem>("Item",   ItemSchema)