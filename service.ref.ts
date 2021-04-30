

import * as bcrypt from "bcrypt";
import { CreateHrDto } from "../dtos/hr.dto";
import HttpException from "../exceptions/HttpException";
import { Hr, HrPopulated } from "../models/hr.models";
import HrModel from "../models/hr.models";
import { IJob } from "../models/jobs.model";
class HrService {
    public hr = HrModel

    public create = async (hrData: CreateHrDto): Promise<Hr> => {
        const findHr: Hr = await this.hr.findOne({ email: hrData.email }) as Hr;
        if (findHr != null) {
            throw new HttpException(409, "Hr with same email already exists");
        }
        hrData.password = await bcrypt.hash(hrData.password, 10);
        const createdHrData = await this.hr.create({ ...hrData, Jobs: [] }) as Hr;
        return createdHrData;
    }

    public find = async (objectId: string): Promise<Hr> => {
        const findHr: Hr = await this.hr.findById(objectId) as Hr;
        if (findHr == null) {
            throw new HttpException(404, "Invalid Hr Id")
        }
        return findHr;
    }

    public findJobs = async (objectId: string, page: number, pageSize: number): Promise<IJob[]> => {
        const findHr: HrPopulated = await this.hr.findById(objectId).populate([{
            path: "Jobs", options: {
                sort: { createdAt: "asc" },
                skip: pageSize * page,
                limit: pageSize
            }
        }]) as HrPopulated;
        if (findHr == null) {
            throw new HttpException(404, "Invalid Hr Id");
        }
        return findHr.Jobs
    }

    public addJob = async (_hrId: string, _jobId: string): Promise<Hr> => {
        await this.hr.update(
            { _id: _hrId },
            {
                $push: { Jobs: _jobId }
            })
        return null;
    }
}   

