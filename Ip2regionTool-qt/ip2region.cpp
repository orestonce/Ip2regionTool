#include "ip2region.h"
#include "ip2region-impl.h"

std::string ConvertDbToTxt(std::string in0, std::string in1){
	std::string in;
	{
		uint32_t length = in0.length();
		char tmp[4];
		tmp[0] = (uint32_t(length) >> 24) & 0xFF;
		tmp[1] = (uint32_t(length) >> 16) & 0xFF;
		tmp[2] = (uint32_t(length) >> 8) & 0xFF;
		tmp[3] = (uint32_t(length) >> 0) & 0xFF;
		in.append(tmp, 4);
		in.append(in0);
	}
	{
		uint32_t length = in1.length();
		char tmp[4];
		tmp[0] = (uint32_t(length) >> 24) & 0xFF;
		tmp[1] = (uint32_t(length) >> 16) & 0xFF;
		tmp[2] = (uint32_t(length) >> 8) & 0xFF;
		tmp[3] = (uint32_t(length) >> 0) & 0xFF;
		in.append(tmp, 4);
		in.append(in1);
	}
	char *out = NULL;
	int outLen = 0;
	Go2cppFn_ConvertDbToTxt((char *)in.data(), in.length(), &out, &outLen);
	std::string retValue;
	int outIdx = 0;
	{
		uint32_t length = 0;
		uint32_t a = uint32_t(uint8_t(out[outIdx+0]) << 24);
		uint32_t b = uint32_t(uint8_t(out[outIdx+1]) << 16);
		uint32_t c = uint32_t(uint8_t(out[outIdx+2]) << 8);
		uint32_t d = uint32_t(uint8_t(out[outIdx+3]) << 0);
		length = a | b | c | d;
		outIdx+=4;
		retValue = std::string(out+outIdx, out+outIdx+length);
		outIdx+=length;
	}
	if (out != NULL) {
		free(out);
	}
	return retValue;
}

std::string ConvertTxtToDb(std::string in0, std::string in1){
	std::string in;
	{
		uint32_t length = in0.length();
		char tmp[4];
		tmp[0] = (uint32_t(length) >> 24) & 0xFF;
		tmp[1] = (uint32_t(length) >> 16) & 0xFF;
		tmp[2] = (uint32_t(length) >> 8) & 0xFF;
		tmp[3] = (uint32_t(length) >> 0) & 0xFF;
		in.append(tmp, 4);
		in.append(in0);
	}
	{
		uint32_t length = in1.length();
		char tmp[4];
		tmp[0] = (uint32_t(length) >> 24) & 0xFF;
		tmp[1] = (uint32_t(length) >> 16) & 0xFF;
		tmp[2] = (uint32_t(length) >> 8) & 0xFF;
		tmp[3] = (uint32_t(length) >> 0) & 0xFF;
		in.append(tmp, 4);
		in.append(in1);
	}
	char *out = NULL;
	int outLen = 0;
	Go2cppFn_ConvertTxtToDb((char *)in.data(), in.length(), &out, &outLen);
	std::string retValue;
	int outIdx = 0;
	{
		uint32_t length = 0;
		uint32_t a = uint32_t(uint8_t(out[outIdx+0]) << 24);
		uint32_t b = uint32_t(uint8_t(out[outIdx+1]) << 16);
		uint32_t c = uint32_t(uint8_t(out[outIdx+2]) << 8);
		uint32_t d = uint32_t(uint8_t(out[outIdx+3]) << 0);
		length = a | b | c | d;
		outIdx+=4;
		retValue = std::string(out+outIdx, out+outIdx+length);
		outIdx+=length;
	}
	if (out != NULL) {
		free(out);
	}
	return retValue;
}

