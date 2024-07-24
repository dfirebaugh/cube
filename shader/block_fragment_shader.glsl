#version 330 core
out vec4 FragColor;

in vec2 TexCoord;
in vec3 Normal;

uniform sampler2D textures[6];

void main() {
    vec3 absNormal = abs(Normal);
    int faceIndex = 0;
    if (absNormal.x > absNormal.y && absNormal.x > absNormal.z) {
        faceIndex = (Normal.x > 0) ? 0 : 1; // Right or Left
    } else if (absNormal.y > absNormal.x && absNormal.y > absNormal.z) {
        faceIndex = (Normal.y > 0) ? 4 : 5; // Top or Bottom
    } else {
        faceIndex = (Normal.z > 0) ? 2 : 3; // Front or Back
    }
    FragColor = texture(textures[faceIndex], TexCoord);
}

