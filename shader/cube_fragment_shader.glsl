#version 330 core
in vec3 fragColor;
out vec4 outputColor;
uniform vec3 inputColor;
void main() {
    outputColor = vec4(inputColor, 1.0);
}

