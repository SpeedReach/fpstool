import tensorrt as trt
import pycuda.driver as cuda
import pycuda.autoinit
import numpy as np

# Initialize TensorRT logger
TRT_LOGGER = trt.Logger(trt.Logger.WARNING)

def load_engine(engine_path):
    with open(engine_path, "rb") as f, trt.Runtime(TRT_LOGGER) as runtime:
        return runtime.deserialize_cuda_engine(f.read())

def allocate_buffers(engine):
    inputs = []
    outputs = []
    bindings = []
    stream = cuda.Stream()

    for binding in engine:
        size = trt.volume(engine.get_tensor_shape(binding)) * 1
        dtype = trt.nptype(engine.get_tensor_dtype(binding))
        host_mem = cuda.pagelocked_empty(size, dtype)
        device_mem = cuda.mem_alloc(host_mem.nbytes)
        bindings.append(int(device_mem))
        if engine.get_tensor_mode(binding) == trt.TensorIOMode.INPUT:
            inputs.append((host_mem, device_mem))
        else:
            outputs.append((host_mem, device_mem))
    return inputs, outputs, bindings, stream

def do_inference(context, bindings, inputs, outputs, stream):
    [cuda.memcpy_htod_async(inp[1], inp[0], stream) for inp in inputs]
    context.execute_async_v3(stream_handle=stream.handle)
    [cuda.memcpy_dtoh_async(out[0], out[1], stream) for out in outputs]
    stream.synchronize()
    return [out[0] for out in outputs]

# Load the TensorRT engine
engine = load_engine("best.engine")

# Allocate buffers
inputs, outputs, bindings, stream = allocate_buffers(engine)

# Create execution context
context = engine.create_execution_context()

# Example input (replace with your actual input)
input_data = np.random.random(size=(1, 3, 416, 416)).astype(np.float32)
np.copyto(inputs[0][0], input_data.ravel())

# Run inference
output_data = do_inference(context, bindings, inputs, outputs, stream)

# Process output_data as needed
print("Inference output:", output_data)
