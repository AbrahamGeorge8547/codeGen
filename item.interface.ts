import {Document} from 'mongoose';


export interface IItem extends Document {
	email: string, 
	price: number, 
	prices: number[], 
	referenceUrl: boolean, 
	website: string
	}

