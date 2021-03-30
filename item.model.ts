import * as mongoose from 'mongoose';


const ItemSchema = new mongoose.Schema({ 
email:{required:true,type:String}, 
price:{required:true,type:Number}, 
prices:{required:true,type:Array}, 
referenceUrl:{required:false,type:Boolean}, 
website:{required:true,type:String} 
 })

export interface IItem extends mongoose.Document {
email:string,
price:number,
prices:number[],
referenceUrl:boolean,
website:string
}

export default mongoose.model<IItem>("Item", ItemSchema)