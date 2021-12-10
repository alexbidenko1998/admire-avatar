import $axios from '~/api/index';
import {ImageInputType, ImageType} from '~/types/image';

export const getPrints = () => $axios.get<ImageType[]>('prints');

export const generatePrints = (data: { count: number } & ImageInputType) => $axios.post('prints', data);