import * as mongoose from "mongoose";
import { IJob } from "./jobs.model";

const CompanySchema = new mongoose.Schema({
    Jobs: [{ type: mongoose.Schema.Types.ObjectId, ref: "Job" }],
    createdAt: Date,
    description: String,
    imageBig: String,
    imageSmall: String,
    tags: Array,
    name: String,
    updatedAt: Date,
    employees: String,
});

interface IBaseCompany extends mongoose.Document {
    name: string;
    imageBig: string;
    imageSmall: string;
    description?: string;
    createdAt: Date;
    updatedAt: Date;
    tags: string[];
    employees: string;
}
export interface ICompany extends IBaseCompany {
    Jobs: string[];
}

export interface IPopulated extends IBaseCompany {
    Jobs: IJob[];
}

const CompanyModel = mongoose.model<ICompany>("Company", CompanySchema);
export default CompanyModel;
