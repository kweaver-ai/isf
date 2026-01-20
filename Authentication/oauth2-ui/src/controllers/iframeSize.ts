import { Request, Response } from "express";

const getIFrameSize = async (_: Request, res: Response) => {
  const height = 410;

  try {
    res.json({
      height,
    });
  } catch (error: any) {
    if (
      error &&
      error?.response &&
      error?.response?.status &&
      error?.response?.data
    ) {
      const { status, data } = error.response;
      res.statusCode = status;
      res.json({ ...data });
    } else {
      res.statusCode = 500;
      res.json({
        code: 500,
        cause: "oauth2-ui内部错误",
        message: "内部错误",
      });
    }
  }
};
export default getIFrameSize;
